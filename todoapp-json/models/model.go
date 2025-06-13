package models

type Status string

// status constants of type Status
const (
	NotStarted Status = "not started"
	Started    Status = "started"
	Completed  Status = "completed"
)

// TodoItem struct type
type TodoItem struct {
	Id          int    `json:"id"`
	Description string `json:"description"`
	Status      Status `json:"status"`
}

func IsValidStatus(status Status) bool {
	switch status {
	case NotStarted, Started, Completed:
		return true
	default:
		return false
	}
}
