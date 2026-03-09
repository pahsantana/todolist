package domain

import "errors"

var (
	TaskNotFound    = errors.New("task not found")
	TaskCompleted   = errors.New("completed tasks cannot be edited")
	InvalidStatus   = errors.New("invalid status")
	InvalidPriority = errors.New("invalid priority")
	DueDateInPast   = errors.New("due_date cannot be in the past")
)
