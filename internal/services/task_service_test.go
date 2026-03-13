package services_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/pahsantana/todolist/internal/domain/entities"
	"github.com/pahsantana/todolist/internal/domain/repository"
	"github.com/pahsantana/todolist/internal/dto"
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

func (m *mockRepo) ListCountByStatus(ctx context.Context) (*dto.TaskSummary, error) {
	summary := &dto.TaskSummary{}
	for _, t := range m.tasks {
		switch t.Status {
		case entities.Pending:
			summary.Pending++
		case entities.InProgress:
			summary.InProgress++
		case entities.Completed:
			summary.Completed++
		case entities.Cancelled:
			summary.Cancelled++
		}
	}
	return summary, nil
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
	task, err := svc.Create(context.Background(), dto.CreateTaskInput{
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
		input   dto.CreateTaskInput
		wantErr error
	}{
		{
			name: "successfully creates a task",
			input: dto.CreateTaskInput{
				Title:    "Study Golang",
				Priority: entities.High,
				DueDate:  futureDate(),
			},
		},
		{
			name:    "fails with invalid priority",
			input:   dto.CreateTaskInput{Title: "Test", Priority: "urgent"},
			wantErr: entities.InvalidPriority,
		},
		{
			name:    "fails when due date is in the past",
			input:   dto.CreateTaskInput{Title: "Test", Priority: entities.Low, DueDate: &past},
			wantErr: entities.DueDateInPast,
		},
		{
			name:    "fails with invalid date format",
			input:   dto.CreateTaskInput{Title: "Test", Priority: entities.Low, DueDate: &invalidDate},
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
		input   dto.UpdateTaskInput
		wantErr error
	}{
		{
			name: "successfully updates a task",
			id: func(svc *services.TaskService, _ *mockRepo) string {
				s := time.Now().AddDate(0, 0, 7).Format(entities.DateLayout)
				task, _ := svc.Create(context.Background(), dto.CreateTaskInput{
					Title:    "Study Golang",
					Priority: entities.High,
					DueDate:  &s,
				})
				return task.ID
			},
			input: dto.UpdateTaskInput{Title: &newTitle, Status: &status},
		},
		{
			name:    "fails when task does not exist",
			id:      func(_ *services.TaskService, _ *mockRepo) string { return "nonexistent-id" },
			input:   dto.UpdateTaskInput{Title: &newTitle},
			wantErr: entities.TaskNotFound,
		},
		{
			name:    "fails when task is completed",
			id:      func(_ *services.TaskService, repo *mockRepo) string { return completedTask(repo).ID },
			input:   dto.UpdateTaskInput{Title: &newTitle},
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

func TestSummary(t *testing.T) {
	tests := []struct {
		name string
		seed func(svc *services.TaskService, repo *mockRepo)
		want dto.TaskSummary
	}{
		{
			name: "returns zero counts when no tasks exist",
			seed: func(_ *services.TaskService, _ *mockRepo) {},
			want: dto.TaskSummary{},
		},
		{
			name: "correctly counts tasks by status",
			seed: func(svc *services.TaskService, repo *mockRepo) {
				s := time.Now().AddDate(0, 0, 7).Format(entities.DateLayout)
				svc.Create(context.Background(), dto.CreateTaskInput{Title: "Task 1", Priority: entities.High, DueDate: &s})
				svc.Create(context.Background(), dto.CreateTaskInput{Title: "Task 2", Priority: entities.Low, DueDate: &s})
				completedTask(repo)
			},
			want: dto.TaskSummary{Pending: 2, Completed: 1},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			svc, repo := newSvc()
			tt.seed(svc, repo)

			summary, err := svc.Summary(context.Background())
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if summary.Pending != tt.want.Pending {
				t.Errorf("pending: got %d, want %d", summary.Pending, tt.want.Pending)
			}
			if summary.InProgress != tt.want.InProgress {
				t.Errorf("in_progress: got %d, want %d", summary.InProgress, tt.want.InProgress)
			}
			if summary.Completed != tt.want.Completed {
				t.Errorf("completed: got %d, want %d", summary.Completed, tt.want.Completed)
			}
			if summary.Cancelled != tt.want.Cancelled {
				t.Errorf("cancelled: got %d, want %d", summary.Cancelled, tt.want.Cancelled)
			}
		})
	}
}
