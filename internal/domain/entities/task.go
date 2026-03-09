package domain

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
