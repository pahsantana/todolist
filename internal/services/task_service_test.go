package services_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/pahsantana/todolist/internal/domain/entities"
	"github.com/pahsantana/todolist/internal/domain/repository"
	"github.com/pahsantana/todolist/internal/services"
)

type mockRepo struct {
	tasks map[string]*entities.Task
}

func newMockRepo() *mockRepo {
	return &mockRepo{tasks: make(map[string]*entities.Task)}
}

func (m *mockRepo) Create(ctx context.Context, task *entities.Task) error {
	m.tasks[task.ID] = task
	return nil
}

func (m *mockRepo) FindAll(ctx context.Context, filters map[string]string) ([]entities.Task, error) {
	var result []entities.Task
	for _, t := range m.tasks {
		result = append(result, *t)
	}
	return result, nil
}

func (m *mockRepo) FindByID(ctx context.Context, id string) (*entities.Task, error) {
	task, ok := m.tasks[id]
	if !ok {
		return nil, nil
	}
	return task, nil
}

func (m *mockRepo) Update(ctx context.Context, id string, task *entities.Task) error {
	if _, ok := m.tasks[id]; !ok {
		return errors.New("task not found")
	}
	m.tasks[id] = task
	return nil
}

func (m *mockRepo) Delete(ctx context.Context, id string) error {
	if _, ok := m.tasks[id]; !ok {
		return errors.New("task not found")
	}
	delete(m.tasks, id)
	return nil
}

var _ repository.TaskRepository = (*mockRepo)(nil)

func newSvc() (*services.TaskService, *mockRepo) {
	repo := newMockRepo()
	return services.NewTaskService(repo), repo
}

func futureDate() *string {
	s := time.Now().AddDate(0, 0, 7).Format(entities.DateLayout)
	return &s
}

func seedTask(t *testing.T, svc *services.TaskService) *entities.Task {
	t.Helper()
	task, err := svc.Create(context.Background(), services.CreateTaskInput{
		Title:    "Study Golang",
		Priority: entities.High,
		DueDate:  futureDate(),
	})
	if err != nil {
		t.Fatalf("failed to seed task: %v", err)
	}
	return task
}

func completedTask(repo *mockRepo) *entities.Task {
	now := time.Now()
	task := &entities.Task{
		ID:        "completed-id",
		Title:     "Done task",
		Status:    entities.Completed,
		Priority:  entities.Low,
		CreatedAt: now,
		UpdatedAt: now,
	}
	repo.tasks[task.ID] = task
	return task
}

func TestCreateTask(t *testing.T) {
	past := "2020-01-01"
	invalidDate := "31-12-2026"

	tests := []struct {
		name    string
		input   services.CreateTaskInput
		wantErr error
	}{
		{
			name: "successfully creates a task",
			input: services.CreateTaskInput{
				Title:    "Study Golang",
				Priority: entities.High,
				DueDate:  futureDate(),
			},
		},
		{
			name:    "fails with invalid priority",
			input:   services.CreateTaskInput{Title: "Test", Priority: "urgent"},
			wantErr: entities.InvalidPriority,
		},
		{
			name:    "fails when due date is in the past",
			input:   services.CreateTaskInput{Title: "Test", Priority: entities.Low, DueDate: &past},
			wantErr: entities.DueDateInPast,
		},
		{
			name:    "fails with invalid date format",
			input:   services.CreateTaskInput{Title: "Test", Priority: entities.Low, DueDate: &invalidDate},
			wantErr: entities.InvalidDateFormat,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			svc, _ := newSvc()
			task, err := svc.Create(context.Background(), tt.input)

			if tt.wantErr != nil {
				if !errors.Is(err, tt.wantErr) {
					t.Errorf("got error %v, want %v", err, tt.wantErr)
				}
				return
			}
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if task.ID == "" {
				t.Error("expected task to have an ID")
			}
			if task.Status != entities.Pending {
				t.Errorf("expected status pending, got %s", task.Status)
			}
		})
	}
}

func TestListTasks(t *testing.T) {
	tests := []struct {
		name    string
		filters map[string]string
		wantErr error
	}{
		{
			name:    "successfully lists all tasks",
			filters: map[string]string{},
		},
		{
			name:    "fails with invalid status filter",
			filters: map[string]string{"status": "invalid"},
			wantErr: entities.InvalidStatus,
		},
		{
			name:    "fails with invalid priority filter",
			filters: map[string]string{"priority": "urgent"},
			wantErr: entities.InvalidPriority,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			svc, _ := newSvc()
			seedTask(t, svc)

			_, err := svc.List(context.Background(), tt.filters)
			if tt.wantErr != nil {
				if !errors.Is(err, tt.wantErr) {
					t.Errorf("got error %v, want %v", err, tt.wantErr)
				}
				return
			}
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
		})
	}
}

func TestGetTaskByID(t *testing.T) {
	tests := []struct {
		name    string
		id      func(created *entities.Task) string
		wantErr error
	}{
		{
			name: "successfully finds a task",
			id:   func(created *entities.Task) string { return created.ID },
		},
		{
			name:    "fails when task does not exist",
			id:      func(_ *entities.Task) string { return "nonexistent-id" },
			wantErr: entities.TaskNotFound,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			svc, _ := newSvc()
			created := seedTask(t, svc)

			_, err := svc.GetByID(context.Background(), tt.id(created))
			if tt.wantErr != nil {
				if !errors.Is(err, tt.wantErr) {
					t.Errorf("got error %v, want %v", err, tt.wantErr)
				}
				return
			}
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
		})
	}
}

func TestUpdateTask(t *testing.T) {
	newTitle := "Updated title"
	status := entities.InProgress

	tests := []struct {
		name    string
		id      func(svc *services.TaskService, repo *mockRepo) string
		input   services.UpdateTaskInput
		wantErr error
	}{
		{
			name: "successfully updates a task",
			id: func(svc *services.TaskService, _ *mockRepo) string {
				s := time.Now().AddDate(0, 0, 7).Format(entities.DateLayout)
				task, _ := svc.Create(context.Background(), services.CreateTaskInput{
					Title:    "Study Golang",
					Priority: entities.High,
					DueDate:  &s,
				})
				return task.ID
			},
			input: services.UpdateTaskInput{Title: &newTitle, Status: &status},
		},
		{
			name:    "fails when task does not exist",
			id:      func(_ *services.TaskService, _ *mockRepo) string { return "nonexistent-id" },
			input:   services.UpdateTaskInput{Title: &newTitle},
			wantErr: entities.TaskNotFound,
		},
		{
			name:    "fails when task is completed",
			id:      func(_ *services.TaskService, repo *mockRepo) string { return completedTask(repo).ID },
			input:   services.UpdateTaskInput{Title: &newTitle},
			wantErr: entities.TaskCompleted,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			svc, repo := newSvc()
			id := tt.id(svc, repo)

			_, err := svc.Update(context.Background(), id, tt.input)
			if tt.wantErr != nil {
				if !errors.Is(err, tt.wantErr) {
					t.Errorf("got error %v, want %v", err, tt.wantErr)
				}
				return
			}
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
		})
	}
}

func TestDeleteTask(t *testing.T) {
	tests := []struct {
		name    string
		id      func(created *entities.Task) string
		wantErr error
	}{
		{
			name: "successfully deletes a task",
			id:   func(created *entities.Task) string { return created.ID },
		},
		{
			name:    "fails when task does not exist",
			id:      func(_ *entities.Task) string { return "nonexistent-id" },
			wantErr: entities.TaskNotFound,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			svc, _ := newSvc()
			created := seedTask(t, svc)

			err := svc.Delete(context.Background(), tt.id(created))
			if tt.wantErr != nil {
				if !errors.Is(err, tt.wantErr) {
					t.Errorf("got error %v, want %v", err, tt.wantErr)
				}
				return
			}
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
		})
	}
}
