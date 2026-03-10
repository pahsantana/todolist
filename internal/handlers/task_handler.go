package handlers

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/pahsantana/todolist/internal/domain/entities"
	"github.com/pahsantana/todolist/internal/services"
	"go.uber.org/zap"
)

const (
	paramID       = "id"
	queryStatus   = "status"
	queryPriority = "priority"
	logTaskID     = "task_id"
	logCreated    = "task created"
	logUpdated    = "task updated"
	logDeleted    = "task deleted"
	logInternal   = "internal error"
)

type TaskHandler struct {
	service *services.TaskService
	log     *zap.Logger
}

func NewTaskHandler(service *services.TaskService, log *zap.Logger) *TaskHandler {
	return &TaskHandler{service: service, log: log}
}

func (h *TaskHandler) Create(c *gin.Context) {
	var input services.CreateTaskInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, errorResponse(err.Error()))
		return
	}

	task, err := h.service.Create(c.Request.Context(), input)
	if err != nil {
		h.handleError(c, err)
		return
	}

	h.log.Info(logCreated, zap.String(logTaskID, task.ID))
	c.JSON(http.StatusCreated, task)
}

func (h *TaskHandler) List(c *gin.Context) {
	filters := map[string]string{
		queryStatus:   c.Query(queryStatus),
		queryPriority: c.Query(queryPriority),
	}

	tasks, err := h.service.List(c.Request.Context(), filters)
	if err != nil {
		h.handleError(c, err)
		return
	}

	c.JSON(http.StatusOK, tasks)
}

func (h *TaskHandler) GetByID(c *gin.Context) {
	task, err := h.service.GetByID(c.Request.Context(), c.Param(paramID))
	if err != nil {
		h.handleError(c, err)
		return
	}

	c.JSON(http.StatusOK, task)
}

func (h *TaskHandler) Update(c *gin.Context) {
	var input services.UpdateTaskInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, errorResponse(err.Error()))
		return
	}

	task, err := h.service.Update(c.Request.Context(), c.Param(paramID), input)
	if err != nil {
		h.handleError(c, err)
		return
	}

	h.log.Info(logUpdated, zap.String(logTaskID, task.ID))
	c.JSON(http.StatusOK, task)
}

func (h *TaskHandler) Delete(c *gin.Context) {
	id := c.Param(paramID)
	if err := h.service.Delete(c.Request.Context(), id); err != nil {
		h.handleError(c, err)
		return
	}

	h.log.Info(logDeleted, zap.String(logTaskID, id))
	c.Status(http.StatusNoContent)
}

func (h *TaskHandler) handleError(c *gin.Context, err error) {
	switch {
	case errors.Is(err, entities.TaskNotFound):
		c.JSON(http.StatusNotFound, errorResponse(err.Error()))
	case errors.Is(err, entities.TaskCompleted):
		c.JSON(http.StatusUnprocessableEntity, errorResponse(err.Error()))
	case errors.Is(err, entities.InvalidStatus),
		errors.Is(err, entities.InvalidPriority),
		errors.Is(err, entities.DueDateInPast):
		c.JSON(http.StatusBadRequest, errorResponse(err.Error()))
	default:
		h.log.Error(logInternal, zap.Error(err))
		c.JSON(http.StatusInternalServerError, errorResponse(entities.InternalServer.Error()))
	}
}

func errorResponse(msg string) gin.H {
	return gin.H{"error": msg}
}