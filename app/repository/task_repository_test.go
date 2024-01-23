package repository_test

import (
	"testing"
	"togo/domain"
	"togo/utils"

	"github.com/stretchr/testify/require"
)

func createRandomTaskManager(t *testing.T) domain.TaskManager {
	user := createRandomUser(t)
	taskManager, rerr := taskRepository.CreateTaskManagerIfNotExists(user.ID)
	require.Empty(t, rerr)
	require.Equal(t, 5, taskManager.TaskLimitPerDay)
	require.Equal(t, user.ID, taskManager.UserID)

	return taskManager
}

func mockTask() domain.Task {
	return domain.Task{
		Name:   utils.RandomName(),
		Status: domain.TaskStatus(utils.RandomInt(0, 2)),
	}
}

func requireEqualTask(t *testing.T, t1 *domain.Task, t2 *domain.Task) {
	require.Equal(t, t1.Name, t2.Name)
	require.Equal(t, t1.Status, t2.Status)
}

func createRandomTask(t *testing.T) domain.Task {
	taskManager := createRandomTaskManager(t)
	randomTask := mockTask()
	task, err := taskRepository.CreateTask(taskManager.UserID, randomTask)
	require.Empty(t, err)
	requireEqualTask(t, &randomTask, &task)
	require.Equal(t, taskManager.ID, task.TaskManagerID)

	return task
}

func TestCreateTaskManager(t *testing.T) {
	createRandomTaskManager(t)
}

func TestGetTaskManagerByUserID(t *testing.T) {
	randomTaskManager := createRandomTaskManager(t)

	taskManager, err := taskRepository.GetTaskManagerByUserID(randomTaskManager.UserID)
	require.Empty(t, err)
	require.Equal(t, randomTaskManager.TaskLimitPerDay, taskManager.TaskLimitPerDay)
	require.Equal(t, randomTaskManager.UserID, taskManager.UserID)
}

// func TestCheckIfTaskManagerExistsByUserID(t *testing.T) {
// 	taskManager := createRandomTaskManager(t)
// 	isExists, err := taskRepository.CheckIfTaskManagerExistsByUserID(taskManager.UserID)
// 	require.Empty(t, err)
// 	require.True(t, isExists)
// }

func TestCreateTaskManagerIfNotExists(t *testing.T) {
	user := createRandomUser(t)

	taskManager1, rerr := taskRepository.CreateTaskManagerIfNotExists(user.ID)
	require.Empty(t, rerr)
	require.NotZero(t, taskManager1.ID)
	require.Equal(t, taskManager1.TaskLimitPerDay, 5)
	require.Equal(t, user.ID, taskManager1.UserID)

	taskManager2, rerr := taskRepository.CreateTaskManagerIfNotExists(user.ID)
	require.Empty(t, rerr)
	require.Equal(t, taskManager1.ID, taskManager2.ID)
	require.Equal(t, taskManager1.TaskLimitPerDay, taskManager2.TaskLimitPerDay)
	require.Equal(t, taskManager1.UserID, taskManager2.UserID)
}

func TestSetTaskLimit(t *testing.T) {
	randomTaskManager := createRandomTaskManager(t)
	randomTaskLimitPerDay := utils.RandomInt(5, 10)

	taskManager, rerr := taskRepository.SetTaskLimit(randomTaskManager.UserID, randomTaskLimitPerDay)
	require.Empty(t, rerr)
	require.Equal(t, randomTaskLimitPerDay, taskManager.TaskLimitPerDay)
	require.Equal(t, randomTaskManager.UserID, taskManager.UserID)
}

func TestCreateTask(t *testing.T) {
	createRandomTask(t)
}

func TestCreateMultiTask(t *testing.T) {
	randomTaskManager := createRandomTaskManager(t)

	for i := 0; i < utils.RandomInt(2, randomTaskManager.TaskLimitPerDay); i++ {
		randomTask := mockTask()
		task, rerr := taskRepository.CreateTask(randomTaskManager.UserID, randomTask)
		require.Empty(t, rerr)
		requireEqualTask(t, &randomTask, &task)
		require.Equal(t, randomTaskManager.ID, task.TaskManagerID)
	}

}

func TestGetTaskByID(t *testing.T) {
	randomTask := createRandomTask(t)

	task, rerr := taskRepository.GetTaskByID(randomTask.ID)
	require.Empty(t, rerr)
	requireEqualTask(t, &randomTask, &task)
	require.Equal(t, randomTask.TaskManagerID, task.TaskManagerID)
}

func TestFetchtaskByUserID(t *testing.T) {
	randomTaskManager := createRandomTaskManager(t)
	var randomTasks []domain.Task
	for i := 0; i < randomTaskManager.TaskLimitPerDay; i++ {
		randomTask := mockTask()
		randomTasks = append(randomTasks, randomTask)
		task, rerr := taskRepository.CreateTask(randomTaskManager.UserID, randomTask)
		require.Empty(t, rerr)
		requireEqualTask(t, &randomTask, &task)
		require.Equal(t, randomTaskManager.ID, task.TaskManagerID)
	}

	tasks, rerr := taskRepository.FetchTaskByUserID(randomTaskManager.UserID)
	require.Empty(t, rerr)
	for i := 0; i < randomTaskManager.TaskLimitPerDay; i++ {
		requireEqualTask(t, &randomTasks[i], &tasks[i])
	}
}

func TestDeleteTask(t *testing.T) {
	randomTask := createRandomTask(t)

	rerr := taskRepository.DeleteTask(randomTask.ID)
	require.Empty(t, rerr)

	task, rerr := taskRepository.GetTaskByID(randomTask.ID)
	require.NotEmpty(t, rerr)
	require.Equal(t, domain.ErrorNotFound, rerr.ErrorType())
	require.NotEqual(t, randomTask.ID, task.ID)
}
