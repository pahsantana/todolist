package handlers

import (
	"errors"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/pahsantana/todolist/internal/domain/entities"
	"github.com/pahsantana/todolist/internal/dto"
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

// Summary godoc
// @Summary     Count tasks by status
// @Tags        tasks
// @Produce     json
// @Success     200 {object} entities.TaskSummary
// @Failure     500 {object} map[string]string
// @Router      /tasks/summary [get]
func (h *TaskHandler) Summary(c *gin.Context) {
	summary, err := h.service.Summary(c.Request.Context())
	if err != nil {
		h.log.Error(logInternal, zap.Error(err))
		c.JSON(http.StatusInternalServerError, errorResponse(entities.InternalServer.Error()))
		return
	}

	c.JSON(http.StatusOK, summary)
}

// Create godoc
// @Summary     Create a new task
// @Tags        tasks
// @Accept      json
// @Produce     json
// @Param       task body dto.CreateTaskInput true "Task data"
// @Success     201 {object} entities.Task
// @Failure     400 {object} map[string]string
// @Failure     500 {object} map[string]string
// @Router      /tasks [post]
func (h *TaskHandler) Create(c *gin.Context) {
	var input dto.CreateTaskInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, errorResponse(translateBindingError(err)))
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

// List godoc
// @Summary     List tasks
// @Tags        tasks
// @Produce     json
// @Param       status   query string false "Filter by status"   Enums(pending, in_progress, completed, cancelled)
// @Param       priority query string false "Filter by priority" Enums(low, medium, high)
// @Success     200 {array}  entities.Task
// @Failure     400 {object} map[string]string
// @Failure     500 {object} map[string]string
// @Router      /tasks [get]
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

// GetByID godoc
// @Summary     Get task by ID
// @Tags        tasks
// @Produce     json
// @Param       id path string true "Task ID"
// @Success     200 {object} entities.Task
// @Failure     404 {object} map[string]string
// @Failure     500 {object} map[string]string
// @Router      /tasks/{id} [get]
func (h *TaskHandler) GetByID(c *gin.Context) {
	task, err := h.service.GetByID(c.Request.Context(), c.Param(paramID))
	if err != nil {
		h.handleError(c, err)
		return
	}

	c.JSON(http.StatusOK, task)
}

// Update godoc
// @Summary     Update a task
// @Tags        tasks
// @Accept      json
// @Produce     json
// @Param       id   path string true "Task ID"
// @Param       task body dto.UpdateTaskInput true "Fields to update"
// @Success     200 {object} entities.Task
// @Failure     400 {object} map[string]string
// @Failure     404 {object} map[string]string
// @Failure     422 {object} map[string]string
// @Failure     500 {object} map[string]string
// @Router      /tasks/{id} [put]
func (h *TaskHandler) Update(c *gin.Context) {
	var input dto.UpdateTaskInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, errorResponse(translateBindingError(err)))
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

// Delete godoc
// @Summary     Delete a task
// @Tags        tasks
// @Param       id path string true "Task ID"
// @Success     204
// @Failure     404 {object} map[string]string
// @Failure     500 {object} map[string]string
// @Router      /tasks/{id} [delete]
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
		errors.Is(err, entities.DueDateInPast),
		errors.Is(err, entities.TitleTooShort),
		errors.Is(err, entities.TitleTooLong),
		errors.Is(err, entities.TitleRequired),
		errors.Is(err, entities.PriorityRequired):
		c.JSON(http.StatusBadRequest, errorResponse(err.Error()))
	default:
		h.log.Error(logInternal, zap.Error(err))
		c.JSON(http.StatusInternalServerError, errorResponse(entities.InternalServer.Error()))
	}
}

func errorResponse(msg string) gin.H {
	return gin.H{"error": msg}
}

func translateBindingError(err error) string {
	msg := err.Error()
	switch {
	case strings.Contains(msg, "Title") && strings.Contains(msg, "min"):
		return entities.TitleTooShort.Error()
	case strings.Contains(msg, "Title") && strings.Contains(msg, "max"):
		return entities.TitleTooLong.Error()
	case strings.Contains(msg, "Title") && strings.Contains(msg, "required"):
		return entities.TitleRequired.Error()
	case strings.Contains(msg, "Priority") && strings.Contains(msg, "required"):
		return entities.PriorityRequired.Error()
	default:
		return msg
	}
}
