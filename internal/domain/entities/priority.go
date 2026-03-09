package domain

type Priority string

const (
	Low    Priority = "low"
	Medium Priority = "medium"
	High   Priority = "high"
)

func IsValidPriority(p Priority) bool {
	switch p {
	case Low, Medium, High:
		return true
	}
	return false
}
