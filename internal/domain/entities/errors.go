package entities

import (
	"errors"
	"fmt"
)

const DateLayout = "2006-01-02"

var (
	TaskNotFound      = errors.New("task not found")
	TaskCompleted     = errors.New("completed tasks cannot be edited")
	InvalidStatus     = errors.New("invalid status")
	InvalidPriority   = errors.New("invalid priority")
	DueDateInPast     = errors.New("due_date cannot be in the past")
	InvalidDateFormat = fmt.Errorf("due_date must be in format %s", DateLayout)
	InternalServer    = errors.New("internal server error")
	TitleRequired  = errors.New("title is required")
	TitleTooShort  = errors.New("title must be at least 3 characters")
	TitleTooLong   = errors.New("title must be at most 100 characters")
	PriorityRequired = errors.New("priority is required")
)
