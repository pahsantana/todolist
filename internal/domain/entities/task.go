package entities

import "time"

type Task struct {
	ID          string    `bson:"_id"         json:"id"`
	Title       string    `bson:"title"       json:"title"`
	Description string    `bson:"description" json:"description,omitempty"`
	Status      Status    `bson:"status"      json:"status"`
	Priority    Priority  `bson:"priority"    json:"priority"`
	DueDate     *string   `bson:"due_date"    json:"due_date,omitempty"`
	CreatedAt   time.Time `bson:"created_at"  json:"created_at"`
	UpdatedAt   time.Time `bson:"updated_at"  json:"updated_at"`
}

func (t *Task) IsCompleted() bool {
	return t.Status == Completed
}

func (t *Task) Apply(title, description *string, status *Status, priority *Priority, dueDate *string) error {
	if title != nil {
		t.Title = *title
	}
	if description != nil {
		t.Description = *description
	}
	if status != nil {
		if !IsValidStatus(*status) {
			return InvalidStatus
		}
		t.Status = *status
	}
	if priority != nil {
		if !IsValidPriority(*priority) {
			return InvalidPriority
		}
		t.Priority = *priority
	}
	if dueDate != nil {
		t.DueDate = dueDate
	}
	t.UpdatedAt = time.Now()
	return nil
}