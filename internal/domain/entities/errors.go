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
	InvalidDateFormat = fmt.orf("due_date must be in format %s", DateLayout)
	InternalServer  = errors.New("internal server error")
)
