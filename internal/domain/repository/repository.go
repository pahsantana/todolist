package repository

import (
    "context"
    "github.com/pahsantana/todolist/internal/domain/entities"
)

type TaskRepository interface {
    Create(ctx context.Context, task *entities.Task) error
    FindAll(ctx context.Context, filters map[string]string) ([]entities.Task, error)
    FindByID(ctx context.Context, id string) (*entities.Task, error)
    Update(ctx context.Context, id string, task *entities.Task) error
    Delete(ctx context.Context, id string) error
}