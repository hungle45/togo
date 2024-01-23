package service_test

import (
	"testing"
	"togo/app/service"
	"togo/domain"
	"togo/domain/mock"
	"togo/utils"

	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
)

func createRandomTask() domain.Task {
	return domain.Task{
		Name:   utils.RandomName(),
		Status: domain.TaskStatus(utils.RandomInt(0, 2)),
	}
}

func requireEqualTask(t *testing.T, t1 *domain.Task, t2 *domain.Task) {
	require.Equal(t, t1.Name, t2.Name)
	require.Equal(t, t1.Status, t2.Status)
}

func TestFetchTaskByUserID(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockTaskRepo := mock.NewMockTaskRepository(ctrl)
	mockUserRepo := mock.NewMockUserRepository(ctrl)
	taskService := service.NewTaskService(mockUserRepo, mockTaskRepo)

	t.Run("Test case 1: Success (tasks non empty)", func(t *testing.T) {
		var userID uint = 0
		var tasks []domain.Task
		for i := 0; i < 5; i++ {
			tasks = append(tasks, createRandomTask())
		}

		mockTaskRepo.EXPECT().
			FetchTaskByUserID(userID).
			Return(tasks, nil)

		resTasks, rerr := taskService.FetchTask(userID)
		require.Nil(t, rerr)
		require.Equal(t, len(tasks), len(resTasks))
		for i := 0; i < 5; i++ {
			requireEqualTask(t, &tasks[i], &resTasks[i])
		}
	})

	t.Run("Test case 2: Success (tasks empty)", func(t *testing.T) {
		var userID uint = 0
		var tasks []domain.Task

		mockTaskRepo.EXPECT().
			FetchTaskByUserID(userID).
			Return(tasks, nil)

		resTasks, rerr := taskService.FetchTask(userID)
		require.Nil(t, rerr)
		require.Equal(t, len(tasks), len(resTasks))
	})

	t.Run("Test case 3: Error (internal)", func(t *testing.T) {
		var userID uint = 0
		mockTaskRepo.EXPECT().
			FetchTaskByUserID(userID).
			Return([]domain.Task{}, domain.NewReponseError(domain.ErrorInternal, "error"))

		_, rerr := taskService.FetchTask(userID)
		require.Equal(t, rerr.ErrorType(), domain.ErrorInternal)
	})
}

func TestGetTaskByID(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockTaskRepo := mock.NewMockTaskRepository(ctrl)
	mockUserRepo := mock.NewMockUserRepository(ctrl)
	taskService := service.NewTaskService(mockUserRepo, mockTaskRepo)

	t.Run("Test case 1: Success", func(t *testing.T) {
		var userID uint = 0
		var taskID uint = 1
		task := createRandomTask()

		mockTaskRepo.EXPECT().
			GetTaskByID(taskID).
			Return(task, nil)

		mockTaskRepo.EXPECT().
			GetTaskManagerByID(task.TaskManagerID).
			Return(domain.TaskManager{UserID: userID}, nil)

		resTask, rerr := taskService.GetTaskByID(userID, taskID)
		require.Nil(t, rerr)
		requireEqualTask(t, &task, &resTask)
	})

	t.Run("Test case 2: Error (task not found)", func(t *testing.T) {
		var userID uint = 0
		var taskID uint = 1

		mockTaskRepo.EXPECT().
			GetTaskByID(taskID).
			Return(domain.Task{}, domain.NewReponseError(domain.ErrorNotFound, ""))

		_, rerr := taskService.GetTaskByID(userID, taskID)
		require.Equal(t, rerr.ErrorType(), domain.ErrorNotFound)
	})

	t.Run("Test case 3: Error (user do not have permistion)", func(t *testing.T) {
		var userID uint = 0
		var taskID uint = 1
		task := createRandomTask()

		mockTaskRepo.EXPECT().
			GetTaskByID(taskID).
			Return(task, nil)

		mockTaskRepo.EXPECT().
			GetTaskManagerByID(task.TaskManagerID).
			Return(domain.TaskManager{UserID: userID + 1}, nil)

		_, rerr := taskService.GetTaskByID(userID, taskID)
		require.Equal(t, rerr.ErrorType(), domain.ErrorPermissionDenied)
	})
}

func TestCreateTask(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockTaskRepo := mock.NewMockTaskRepository(ctrl)
	mockUserRepo := mock.NewMockUserRepository(ctrl)
	taskService := service.NewTaskService(mockUserRepo, mockTaskRepo)

	t.Run("Test case 1: Success", func(t *testing.T) {
		var userID uint = 0
		task := createRandomTask()

		mockTaskRepo.EXPECT().
			CreateTask(userID, task).
			Return(task, nil)

		resTask, rerr := taskService.CreateTask(userID, task)
		require.Nil(t, rerr)
		requireEqualTask(t, &task, &resTask)
	})

	t.Run("Test case 2: Error (internal)", func(t *testing.T) {
		var userID uint = 0
		task := createRandomTask()

		mockTaskRepo.EXPECT().
			CreateTask(userID, task).
			Return(task, domain.NewReponseError(domain.ErrorInternal, ""))

		_, rerr := taskService.CreateTask(userID, task)
		require.Equal(t, rerr.ErrorType(), domain.ErrorInternal)
	})
}

func TestUpdateTask(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockTaskRepo := mock.NewMockTaskRepository(ctrl)
	mockUserRepo := mock.NewMockUserRepository(ctrl)
	taskService := service.NewTaskService(mockUserRepo, mockTaskRepo)

	t.Run("Test case 1: Success", func(t *testing.T) {
		var userID uint = 0
		var taskID uint = 1
		task := createRandomTask()

		mockTaskRepo.EXPECT().
			GetTaskByID(taskID).
			Return(task, nil)

		mockTaskRepo.EXPECT().
			GetTaskManagerByID(task.TaskManagerID).
			Return(domain.TaskManager{UserID: userID}, nil)

		mockTaskRepo.EXPECT().
			UpdateTask(taskID, task).
			Return(task, nil)

		resTask, rerr := taskService.UpdateTask(userID, taskID, task)
		require.Nil(t, rerr)
		requireEqualTask(t, &task, &resTask)
	})

	t.Run("Test case 2: Error (task not found)", func(t *testing.T) {
		var userID uint = 0
		var taskID uint = 1
		task := createRandomTask()

		mockTaskRepo.EXPECT().
			GetTaskByID(taskID).
			Return(domain.Task{}, domain.NewReponseError(domain.ErrorNotFound, ""))

		_, rerr := taskService.UpdateTask(userID, taskID, task)
		require.Equal(t, rerr.ErrorType(), domain.ErrorNotFound)
	})

	t.Run("Test case 3: Error (user do not have permistion)", func(t *testing.T) {
		var userID uint = 0
		var taskID uint = 1
		task := createRandomTask()

		mockTaskRepo.EXPECT().
			GetTaskByID(taskID).
			Return(task, nil)

		mockTaskRepo.EXPECT().
			GetTaskManagerByID(task.TaskManagerID).
			Return(domain.TaskManager{UserID: userID + 1}, nil)

		_, rerr := taskService.UpdateTask(userID, taskID, task)
		require.Equal(t, rerr.ErrorType(), domain.ErrorPermissionDenied)
	})
}

func TestDeleteTask(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockTaskRepo := mock.NewMockTaskRepository(ctrl)
	mockUserRepo := mock.NewMockUserRepository(ctrl)
	taskService := service.NewTaskService(mockUserRepo, mockTaskRepo)

	t.Run("Test case 1: Success", func(t *testing.T) {
		var userID uint = 0
		var taskID uint = 1
		task := createRandomTask()

		mockTaskRepo.EXPECT().
			GetTaskByID(taskID).
			Return(task, nil)

		mockTaskRepo.EXPECT().
			GetTaskManagerByID(task.TaskManagerID).
			Return(domain.TaskManager{UserID: userID}, nil)

		mockTaskRepo.EXPECT().
			DeleteTask(taskID).
			Return(nil)

		rerr := taskService.DeleteTask(userID, taskID)
		require.Nil(t, rerr)
	})

	t.Run("Test case 2: Error (internal)", func(t *testing.T) {
		var userID uint = 0
		var taskID uint = 1
		task := createRandomTask()

		mockTaskRepo.EXPECT().
			GetTaskByID(taskID).
			Return(task, nil)

		mockTaskRepo.EXPECT().
			GetTaskManagerByID(task.TaskManagerID).
			Return(domain.TaskManager{UserID: userID}, nil)

		mockTaskRepo.EXPECT().
			DeleteTask(taskID).
			Return(domain.NewReponseError(domain.ErrorInternal, ""))

		rerr := taskService.DeleteTask(userID, taskID)
		require.Equal(t, rerr.ErrorType(), domain.ErrorInternal)
	})
}

func TestUpdateTaskLimit(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockTaskRepo := mock.NewMockTaskRepository(ctrl)
	mockUserRepo := mock.NewMockUserRepository(ctrl)
	taskService := service.NewTaskService(mockUserRepo, mockTaskRepo)

	t.Run("Test case 1: Success", func(t *testing.T) {
		var userID uint = 0
		var targetUserID uint = 1
		var taskLimitPerDay int = 2

		mockUserRepo.EXPECT().
			IsAdmin(userID).
			Return(true, nil)

		mockTaskRepo.EXPECT().
			SetTaskLimit(targetUserID, taskLimitPerDay).
			Return(domain.TaskManager{}, nil)

		rerr := taskService.UpdateTaskLimit(userID, targetUserID, taskLimitPerDay)
		require.Nil(t, rerr)
	})

	t.Run("Test case 2: Error (user is not admin)", func(t *testing.T) {
		var userID uint = 0
		var targetUserID uint = 1
		var taskLimitPerDay int = 2

		mockUserRepo.EXPECT().
			IsAdmin(userID).
			Return(false, nil)

		rerr := taskService.UpdateTaskLimit(userID, targetUserID, taskLimitPerDay)
		require.Equal(t, rerr.ErrorType(), domain.ErrorPermissionDenied)
	})

	t.Run("Test case 3: Error (internal)", func(t *testing.T) {
		var userID uint = 0
		var targetUserID uint = 1
		var taskLimitPerDay int = 2

		mockUserRepo.EXPECT().
			IsAdmin(userID).
			Return(true, nil)

		mockTaskRepo.EXPECT().
			SetTaskLimit(targetUserID, taskLimitPerDay).
			Return(domain.TaskManager{}, domain.NewReponseError(domain.ErrorInternal, ""))

		rerr := taskService.UpdateTaskLimit(userID, targetUserID, taskLimitPerDay)
		require.Equal(t, rerr.ErrorType(), domain.ErrorInternal)
	})
}
