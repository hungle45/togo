package http

import (
	"fmt"
	"net/http"
	"strconv"

	http_utils "togo/app/delivery/http/utils"
	"togo/domain"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
)

type taskForm struct {
	Name   string            `json:"name" binding:"required"`
	Status domain.TaskStatus `json:"status" binding:"required"`
}

type TaskLimitForm struct {
	UserID          uint `json:"user_id" binding:"required"`
	TaskLimitPerDay int  `json:"task_limit_per_day" binding:"required"`
}

type taskHTTPHandler struct {
	taskService domain.TaskService
}

func NewTaskHTTPHandler(r *gin.RouterGroup, taskService domain.TaskService, jwtMiddleware gin.HandlerFunc) {
	handler := taskHTTPHandler{taskService: taskService}
	
	taskRouter := r.Group("/tasks")

	userGroup := taskRouter.Group("", jwtMiddleware)
	{
		userGroup.GET("/", handler.list)
		userGroup.GET("/:taskID", handler.getByID)
		userGroup.POST("/", handler.create)
		userGroup.PUT("/:taskID", handler.update)
		userGroup.DELETE("/:taskID", handler.delete)
	}

	adminGroup := taskRouter.Group("", jwtMiddleware)
	{
		adminGroup.POST("/limit", handler.updateTaskLimit)
	}
}

func (handler *taskHTTPHandler) list(c *gin.Context) {
	userID, err := getUserIDFromContext(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, http_utils.ResponseWithMessage(
			http_utils.ResponseStatusFail, err.Error(),
		))
		return
	}

	tasks, rerr := handler.taskService.FetchTask(userID)
	if rerr != nil {
		c.JSON(http_utils.GetStatusCode(rerr), http_utils.ResponseWithMessage(
			http_utils.ResponseStatusFail, rerr.Message(),
		))
		return
	}

	c.JSON(http.StatusOK, http_utils.ResponseWithData(
		http_utils.ReponseStatusSuccess, map[string]interface{}{
			"tasks": tasks,
		},
	))
}

func (handler *taskHTTPHandler) getByID(c *gin.Context) {
	userID, err := getUserIDFromContext(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, http_utils.ResponseWithMessage(
			http_utils.ResponseStatusFail, err.Error(),
		))
	}

	taskID, err := getTaskIDFromParam(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, http_utils.ResponseWithMessage(
			http_utils.ResponseStatusFail, err.Error(),
		))
		return
	}

	task, rerr := handler.taskService.GetTaskByID(userID, taskID)
	if rerr != nil {
		c.JSON(http_utils.GetStatusCode(rerr), http_utils.ResponseWithMessage(
			http_utils.ResponseStatusFail, rerr.Message(),
		))
		return
	}

	c.JSON(http.StatusOK, http_utils.ResponseWithData(
		http_utils.ReponseStatusSuccess, map[string]interface{}{
			"task": task,
		},
	))
}

func (handler *taskHTTPHandler) create(c *gin.Context) {
	c.Request.Context()
	userID, err := getUserIDFromContext(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, http_utils.ResponseWithMessage(
			http_utils.ResponseStatusFail, err.Error(),
		))
	}

	form, err := validateTaskForm(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, http_utils.ResponseWithMessage(
			http_utils.ResponseStatusFail, err.Error(),
		))
		return
	}

	task := domain.Task{
		Name:   form.Name,
		Status: form.Status,
	}

	_, rerr := handler.taskService.CreateTask(userID, task)
	if rerr != nil {
		c.JSON(http_utils.GetStatusCode(rerr), http_utils.ResponseWithMessage(
			http_utils.ResponseStatusFail, rerr.Message(),
		))
		return
	}

	c.JSON(http.StatusOK, http_utils.ResponseWithMessage(
		http_utils.ReponseStatusSuccess, "task has been created",
	))
}

func (handler *taskHTTPHandler) update(c *gin.Context) {
	userID, err := getUserIDFromContext(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, http_utils.ResponseWithMessage(
			http_utils.ResponseStatusFail, err.Error(),
		))
	}

	taskID, err := getTaskIDFromParam(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, http_utils.ResponseWithMessage(
			http_utils.ResponseStatusFail, err.Error(),
		))
		return
	}

	form, err := validateTaskForm(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, http_utils.ResponseWithMessage(
			http_utils.ResponseStatusFail, err.Error(),
		))
		return
	}

	task := domain.Task{
		Name:   form.Name,
		Status: form.Status,
	}

	_, rerr := handler.taskService.UpdateTask(userID, taskID, task)
	if rerr != nil {
		c.JSON(http_utils.GetStatusCode(rerr), http_utils.ResponseWithMessage(
			http_utils.ResponseStatusFail, rerr.Message(),
		))
		return
	}

	c.JSON(http.StatusOK, http_utils.ResponseWithMessage(
		http_utils.ReponseStatusSuccess, fmt.Sprintf("task with ID %v has been updated", taskID),
	))
}

func (handler *taskHTTPHandler) delete(c *gin.Context) {
	userID, err := getUserIDFromContext(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, http_utils.ResponseWithMessage(
			http_utils.ResponseStatusFail, err.Error(),
		))
	}

	taskID, err := getTaskIDFromParam(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, http_utils.ResponseWithMessage(
			http_utils.ResponseStatusFail, err.Error(),
		))
		return
	}

	rerr := handler.taskService.DeleteTask(userID, taskID)
	if rerr != nil {
		c.JSON(http_utils.GetStatusCode(rerr), http_utils.ResponseWithMessage(
			http_utils.ResponseStatusFail, err.Error(),
		))
		return
	}

	c.JSON(http.StatusOK, http_utils.ResponseWithMessage(
		http_utils.ReponseStatusSuccess, fmt.Sprintf("task with ID %v has been deleted", taskID),
	))
}

func (handler *taskHTTPHandler) updateTaskLimit(c *gin.Context) {
	userID, err := getUserIDFromContext(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, http_utils.ResponseWithMessage(
			http_utils.ResponseStatusFail, err.Error(),
		))
	}

	form, err := validateTaskLimitForm(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, http_utils.ResponseWithMessage(
			http_utils.ResponseStatusFail, err.Error(),
		))
		return
	}

	rerr := handler.taskService.UpdateTaskLimit(userID, form.UserID, form.TaskLimitPerDay)
	if rerr != nil {
		c.JSON(http_utils.GetStatusCode(rerr), http_utils.ResponseWithMessage(
			http_utils.ResponseStatusFail, rerr.Message(),
		))
		return
	}

	c.JSON(http.StatusOK, http_utils.ResponseWithMessage(
		http_utils.ReponseStatusSuccess, fmt.Sprintf("task limit has been updated to %v", form.TaskLimitPerDay),
	))
}

func validateTaskForm(c *gin.Context) (taskForm, error) {
	form := taskForm{}
	if err := c.ShouldBindBodyWith(&form, binding.JSON); err != nil {
		return form, err
	}
	return form, nil
}

func validateTaskLimitForm(c *gin.Context) (TaskLimitForm, error) {
	form := TaskLimitForm{}
	if err := c.ShouldBindBodyWith(&form, binding.JSON); err != nil {
		return form, err
	}
	return form, nil
}

func getUserIDFromContext(c *gin.Context) (uint, error) {
	userID, ok := c.MustGet("userID").(uint)
	if !ok {
		return 0, fmt.Errorf("user ID must be type of uint")
	}
	return userID, nil
}

func getTaskIDFromParam(c *gin.Context) (uint, error) {
	taskID, err := strconv.Atoi(c.Param("taskID"))
	if err != nil {
		return 0, fmt.Errorf("task ID must be a valid integer")
	}
	if taskID < 0 {
		return 0, fmt.Errorf("task ID must be greater than 0")
	}
	return uint(taskID), nil
}
