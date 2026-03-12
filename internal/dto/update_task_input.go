package dto

import (
	"github.com/pahsantana/todolist/internal/domain/entities"
)

type UpdateTaskInput struct {
	Title       *string            `json:"title"       binding:"omitempty,min=3,max=100"`
	Description *string            `json:"description"`
	Status      *entities.Status   `json:"status"`
	Priority    *entities.Priority `json:"priority"`
	DueDate     *string            `json:"due_date"`
}