package usecase

import (
	"fmt"
	"togo/domain"
)

type taskService struct {
	userRepo domain.UserRepository
	taskRepo domain.TaskReponsitory
}

func NewTaskService(userRepo domain.UserRepository, taskRepo domain.TaskReponsitory) domain.TaskService {
	return &taskService{userRepo: userRepo, taskRepo: taskRepo}
}

func (tS *taskService) FetchTask(userID uint) ([]domain.Task, domain.ResponseError) {
	tasks, rerr := tS.taskRepo.FetchTaskByUserID(userID)
	if rerr != nil {
		return tasks, rerr
	}

	return tasks, nil
}

func (tS *taskService) GetTaskByID(userID uint, taskID uint) (domain.Task, domain.ResponseError) {
	task, rerr := tS.taskRepo.GetTaskByID(taskID)
	if rerr != nil && rerr.ErrorType() != domain.ErrorNotFound {
		return domain.Task{}, rerr
	}

	taskManager, rerr := tS.taskRepo.GetTaskManagerByID(task.TaskManagerID)
	if rerr != nil {
		return domain.Task{}, rerr
	}

	if taskManager.UserID != userID {
		return domain.Task{}, domain.NewReponseError(
			domain.ErrorPermissionDenied, fmt.Sprintf("task with id %d not found", taskID),
		)
	}

	return task, nil
}

func (tS *taskService) CreateTask(userID uint, task domain.Task) (domain.Task, domain.ResponseError) {
	// var rerr domain.ResponseError
	// deadline := time.Now().Add(5 * time.Second)
	// for time.Now().Before(deadline) {
	// 	task, rerr = tS.taskRepo.CreateTask(userID, task)
	// 	if rerr == nil {
	// 		return task, nil
	// 	}
	// 	time.Sleep(100 * time.Millisecond)
	// }

	// return domain.Task{}, rerr
	task, rerr := tS.taskRepo.CreateTask(userID, task)
	if rerr != nil {
		return domain.Task{}, rerr
	}

	return task, nil
}

func (tS *taskService) UpdateTask(userID uint, taskID uint, task domain.Task) (domain.Task, domain.ResponseError) {
	_, rerr := tS.GetTaskByID(userID, taskID)
	if rerr != nil {
		return domain.Task{}, rerr
	}

	task, rerr = tS.taskRepo.UpdateTask(taskID, task)
	if rerr != nil {
		return domain.Task{}, rerr
	}

	return task, nil
}

func (tS *taskService) DeleteTask(userID uint, taskID uint) domain.ResponseError {
	rerr := tS.taskRepo.DeleteTask(taskID)
	if rerr != nil {
		return rerr
	}
	return nil
}

func (tS *taskService) UpdateTaskLimit(userID uint, targetUserID uint, taskLimitPerDay int) domain.ResponseError {
	// Check if the user is an admin
	isAdmin, rerr := tS.userRepo.IsAdmin(userID)
	if rerr != nil {
		return rerr
	}

	if !isAdmin {
		return domain.NewReponseError(
			domain.ErrorPermissionDenied, "only admin can update task limit",
		)
	}

	_, rerr = tS.taskRepo.SetTaskLimit(targetUserID, taskLimitPerDay)
	if rerr != nil {
		return rerr
	}

	return nil
}
