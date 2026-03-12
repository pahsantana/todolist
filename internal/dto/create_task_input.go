package dto

import (
	"github.com/pahsantana/todolist/internal/domain/entities"
)

type CreateTaskInput struct {
	Title       string            `json:"title"       binding:"required,min=3,max=100"`
	Description string            `json:"description"`
	Priority    entities.Priority `json:"priority"    binding:"required"`
	DueDate     *string           `json:"due_date"`
}