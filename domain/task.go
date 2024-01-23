package domain

import (
	"gorm.io/gorm"
)

type TaskStatus int

const (
	ToDo       TaskStatus = 1
	Processing TaskStatus = 2
	Done       TaskStatus = 3
)

type TaskManager struct {
	gorm.Model
	TaskLimitPerDay int    `json:"task_limit_per_day"`
	UserID          uint   `json:"user_id" gorm:"index"`
	User            User   `json:"-" gorm:"foreignKey:UserID"`
	Tasks           []Task `json:"task" gorm:"foreignKey:TaskManagerID"`
}

type Task struct {
	gorm.Model
	Name          string     `json:"name"`
	Status        TaskStatus `json:"status"`
	TaskManagerID uint       `json:"task_manager_id" gorm:"index"`
}

type TaskService interface {
	FetchTask(userID uint) (tasks []Task, rerr ResponseError)
	GetTaskByID(userID uint, taskID uint) (res Task, rerr ResponseError)
	CreateTask(userID uint, task Task) (res Task, rerr ResponseError)
	UpdateTask(userID uint, taskID uint, task Task) (res Task, rerr ResponseError)
	DeleteTask(userID uint, taskID uint) (rerr ResponseError)
	UpdateTaskLimit(userID uint, targetUserID uint, taskLimitPerDay int) (rerr ResponseError)
}

type TaskReponsitory interface {
	FetchTaskByUserID(userID uint) (tasks []Task, rerr ResponseError)
	GetTaskByID(taskID uint) (res Task, rerr ResponseError)
	CreateTask(userID uint, task Task) (res Task, rerr ResponseError)
	UpdateTask(taskID uint, task Task) (res Task, rerr ResponseError)
	DeleteTask(taskID uint) (rerr ResponseError)
	SetTaskLimit(userID uint, taskLimitPerDay int) (res TaskManager, rerr ResponseError)
	GetTaskManagerByID(taskManagerID uint) (res TaskManager, rerr ResponseError)
	GetTaskManagerByUserID(userID uint) (res TaskManager, rerr ResponseError)
	CreateTaskManagerIfNotExists(userID uint) (res TaskManager, rerr ResponseError)
}
