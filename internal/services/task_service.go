package services

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/pahsantana/todolist/internal/domain/entities"
	"github.com/pahsantana/todolist/internal/domain/repository"
	"github.com/pahsantana/todolist/internal/dto"
)

type TaskService struct {
	repo repository.TaskRepository
}

func NewTaskService(repo repository.TaskRepository) *TaskService {
	return &TaskService{repo: repo}
}

func (uc *TaskService) Summary(ctx context.Context) (*dto.TaskSummary, error) {
	return uc.repo.ListCountByStatus(ctx)
}

func (uc *TaskService) Create(ctx context.Context, input dto.CreateTaskInput) (*entities.Task, error) {
	if !entities.IsValidPriority(input.Priority) {
		return nil, entities.InvalidPriority
	}
	if input.DueDate != nil {
		if err := validateFutureDate(*input.DueDate); err != nil {
			return nil, err
		}
	}

	now := time.Now()
	task := &entities.Task{
		ID:          uuid.NewString(),
		Title:       input.Title,
		Description: input.Description,
		Status:      entities.Pending,
		Priority:    input.Priority,
		DueDate:     input.DueDate,
		CreatedAt:   now,
		UpdatedAt:   now,
	}

	return task, uc.repo.Create(ctx, task)
}

func (uc *TaskService) List(ctx context.Context, filters map[string]string) ([]entities.Task, error) {
	if status, ok := filters["status"]; ok && status != "" {
		if !entities.IsValidStatus(entities.Status(status)) {
			return nil, entities.InvalidStatus
		}
	}
	if priority, ok := filters["priority"]; ok && priority != "" {
		if !entities.IsValidPriority(entities.Priority(priority)) {
			return nil, entities.InvalidPriority
		}
	}

	tasks, err := uc.repo.FindAll(ctx, filters)
	if err != nil {
		return nil, err
	}
	if tasks == nil {
		tasks = []entities.Task{}
	}
	return tasks, nil
}

func (uc *TaskService) GetByID(ctx context.Context, id string) (*entities.Task, error) {
	task, err := uc.repo.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if task == nil {
		return nil, entities.TaskNotFound
	}
	return task, nil
}

func (uc *TaskService) Update(ctx context.Context, id string, input dto.UpdateTaskInput) (*entities.Task, error) {
	task, err := uc.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if task.IsCompleted() {
		return nil, entities.TaskCompleted
	}
	if input.DueDate != nil {
		if err := validateFutureDate(*input.DueDate); err != nil {
			return nil, err
		}
	}

	if err := task.Apply(input.Title, input.Description, input.Status, input.Priority, input.DueDate); err != nil {
		return nil, err
	}

	return task, uc.repo.Update(ctx, id, task)
}

func (uc *TaskService) Delete(ctx context.Context, id string) error {
	if _, err := uc.GetByID(ctx, id); err != nil {
		return err
	}
	return uc.repo.Delete(ctx, id)
}

func validateFutureDate(dateStr string) error {
    parsed, err := time.Parse(entities.DateLayout, dateStr)
    if err != nil {
        return entities.InvalidDateFormat
    }
    if parsed.Before(time.Now().Truncate(24 * time.Hour)) {
        return entities.DueDateInPast
    }
    return nil
}