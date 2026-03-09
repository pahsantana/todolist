package domain

type Status string

const (
	Pending    Status = "pending"
	InProgress Status = "in_progress"
	Completed  Status = "completed"
	Cancelled  Status = "cancelled"
)

func IsValidStatus(s Status) bool {
	switch s {
	case Pending, InProgress, Completed, Cancelled:
		return true
	}
	return false
}
