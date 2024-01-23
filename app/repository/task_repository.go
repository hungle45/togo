package repository

import (
	"database/sql"
	"errors"
	"os"
	"strconv"
	"time"
	"togo/domain"

	"gorm.io/gorm"
)

type taskRepository struct {
	conn *gorm.DB
}

func NewTaskRepository(conn *gorm.DB) domain.TaskRepository {
	return &taskRepository{conn: conn}
}

func (tRepo *taskRepository) GetTaskByID(taskID uint) (domain.Task, domain.ResponseError) {
	var task domain.Task

	result := tRepo.conn.Where("id = ?", taskID).First(&task)
	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return domain.Task{}, domain.NewReponseError(
			domain.ErrorNotFound, result.Error.Error(),
		)
	}
	if result.Error != nil {
		return domain.Task{}, domain.NewReponseError(
			domain.ErrorInternal, result.Error.Error(),
		)
	}

	return task, nil
}

func (tRepo *taskRepository) FetchTaskByUserID(userID uint) ([]domain.Task, domain.ResponseError) {
	taskManager, rerr := tRepo.CreateTaskManagerIfNotExists(userID)
	if rerr != nil {
		return []domain.Task{}, rerr
	}

	var tasks []domain.Task
	err := tRepo.conn.Model(&taskManager).Association("Tasks").Find(&tasks)
	if err != nil {
		return []domain.Task{}, domain.NewReponseError(
			domain.ErrorInternal, err.Error(),
		)
	}

	return tasks, nil
}

func (tRepo *taskRepository) CreateTask(userID uint, task domain.Task) (domain.Task, domain.ResponseError) {
	taskManager, rerr := tRepo.CreateTaskManagerIfNotExists(userID)
	if rerr != nil {
		return domain.Task{}, rerr
	}

	// tx := tRepo.conn.Begin(&sql.TxOptions{Isolation: sql.LevelSerializable})
	tx := tRepo.conn.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	// create task
	err := tx.Model(&taskManager).Association("Tasks").Append(&task)
	if err != nil {
		tx.Rollback()
		return domain.Task{}, domain.NewReponseError(
			domain.ErrorInternal, err.Error(),
		)
	}

	// check limit task
	now := time.Now().In(tx.Config.NowFunc().Location())
	startTime := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())
	endTime := time.Date(now.Year(), now.Month(), now.Day(), 23, 59, 59, 999999999, now.Location())
	todayTaskCount := tx.Model(&taskManager).Where("created_at BETWEEN ? AND ?", startTime, endTime).Association("Tasks").Count()
	if err != nil {
		tx.Rollback()
		return domain.Task{}, domain.NewReponseError(
			domain.ErrorInternal, err.Error(),
		)
	}

	if int(todayTaskCount) > taskManager.TaskLimitPerDay {
		tx.Rollback()
		return domain.Task{}, domain.NewReponseError(
			domain.ErrorResourceExhausted, "task limit exceeded",
		)
	}

	tx.Commit()
	return task, nil
}

func (tRepo *taskRepository) UpdateTask(taskID uint, task domain.Task) (domain.Task, domain.ResponseError) {
	oldTask, rerr := tRepo.GetTaskByID(taskID)
	if rerr != nil {
		return domain.Task{}, rerr
	}

	tx := tRepo.conn.Model(&oldTask).Updates(map[string]interface{}{
		"name":   task.Name,
		"status": task.Status,
	})
	if tx.Error != nil {
		return domain.Task{}, domain.NewReponseError(
			domain.ErrorInternal, tx.Error.Error(),
		)
	}

	return task, nil
}

func (tRepo *taskRepository) DeleteTask(taskID uint) domain.ResponseError {
	task, rerr := tRepo.GetTaskByID(taskID)
	if rerr != nil {
		if rerr.ErrorType() == domain.ErrorNotFound {
			return nil
		}
		return rerr
	}

	result := tRepo.conn.Unscoped().Delete(&task) // delete pernamently
	if result.Error != nil {
		return domain.NewReponseError(
			domain.ErrorInternal, result.Error.Error(),
		)
	}

	return nil
}

// func (tRepo *taskRepository) CountTodayTaskByUserID(userID uint) (int, domain.ResponseError) {
// 	taskManager, rerr := tRepo.CreateTaskManagerIfNotExists(userID)
// 	if rerr != nil {
// 		return 0, rerr
// 	}

// 	// check limit task
// 	var todayTasks []domain.Task
// 	now := time.Now().In(tRepo.conn.Config.NowFunc().Location())
// 	startTime := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())
// 	endTime := time.Date(now.Year(), now.Month(), now.Day(), 23, 59, 59, 999999999, now.Location())
// 	// todayTaskCount := tRepo.conn.Model(&taskManager).Where("created_at BETWEEN ? AND ?", startTime, endTime).Association("Tasks").Count()
// 	err := tRepo.conn.Model(&taskManager).Where("created_at BETWEEN ? AND ?", startTime, endTime).Association("Tasks").Find(&todayTasks)
// 	if err != nil {
// 		return 0, domain.NewReponseError(
// 			domain.ErrorInternal, err.Error(),
// 		)
// 	}

// 	return len(todayTasks), nil
// }

func (tRepo *taskRepository) SetTaskLimit(userID uint, taskLimitPerDay int) (domain.TaskManager, domain.ResponseError) {
	taskManager, rerr := tRepo.CreateTaskManagerIfNotExists(userID)
	if rerr != nil {
		return domain.TaskManager{}, rerr
	}

	result := tRepo.conn.Model(&taskManager).Update("task_limit_per_day", taskLimitPerDay)
	if result.Error != nil {
		return domain.TaskManager{}, domain.NewReponseError(
			domain.ErrorInternal, result.Error.Error(),
		)
	}

	return taskManager, nil
}

func (tRepo *taskRepository) CreateTaskManagerIfNotExists(userID uint) (domain.TaskManager, domain.ResponseError) {
	var taskManager domain.TaskManager

	tx := tRepo.conn.Begin(&sql.TxOptions{Isolation: sql.LevelSerializable})
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	result := tx.Where("user_id = ?", userID).First(&taskManager)

	defaultTaskLimitPerDay, err := strconv.Atoi(os.Getenv("DEFAULT_TASK_LIMIT_PER_DAY"))
	if err != nil {
		defaultTaskLimitPerDay = 5
	}

	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		taskManager = domain.TaskManager{
			TaskLimitPerDay: defaultTaskLimitPerDay,
			UserID:          userID,
		}
		rerr := tx.Create(&taskManager).Error
		if rerr != nil {
			tx.Rollback()
			return domain.TaskManager{}, domain.NewReponseError(
				domain.ErrorInternal, result.Error.Error(),
			)
		}
	} else if result.Error != nil {
		tx.Rollback()
		return domain.TaskManager{}, domain.NewReponseError(
			domain.ErrorInternal, result.Error.Error(),
		)
	}

	tx.Commit()
	return taskManager, nil
}

func (tRepo *taskRepository) GetTaskManagerByID(taskManagerID uint) (domain.TaskManager, domain.ResponseError) {
	var taskManager domain.TaskManager

	result := tRepo.conn.Where("id = ?", taskManagerID).First(&taskManager)
	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return domain.TaskManager{}, domain.NewReponseError(
			domain.ErrorNotFound, result.Error.Error(),
		)
	}
	if result.Error != nil {
		return domain.TaskManager{}, domain.NewReponseError(
			domain.ErrorInternal, result.Error.Error(),
		)
	}

	return taskManager, nil
}

func (tRepo *taskRepository) GetTaskManagerByUserID(userID uint) (domain.TaskManager, domain.ResponseError) {
	var taskManager domain.TaskManager

	result := tRepo.conn.Where("user_id = ?", userID).First(&taskManager)
	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return domain.TaskManager{}, domain.NewReponseError(
			domain.ErrorNotFound, result.Error.Error(),
		)
	}
	if result.Error != nil {
		return domain.TaskManager{}, domain.NewReponseError(
			domain.ErrorInternal, result.Error.Error(),
		)
	}

	return taskManager, nil
}
