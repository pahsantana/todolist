package domain

import "context"
import "todo-list/internal/domain/entities"

type TaskRepository interface {
	Create(ctx context.Context, task *Task) error
	FindAll(ctx context.Context, filters map[string]string) ([]Task, error)
	FindByID(ctx context.Context, id string) (*Task, error)
	Update(ctx context.Context, id string, task *Task) error
	Delete(ctx context.Context, id string) error
}