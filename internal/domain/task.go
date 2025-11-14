package domain

import "time"

type Task struct {
	ID        string    `json:"id"`
	ListID    string    `json:"list_id"`
	Text      string    `json:"text"`
	Completed bool      `json:"completed"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type CreateTaskRequest struct {
	Text string `json:"text"`
}

type UpdateTaskRequest struct {
	Text      *string `json:"text,omitempty"`
	Completed *bool   `json:"completed,omitempty"`
}
