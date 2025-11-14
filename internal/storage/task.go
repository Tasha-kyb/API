package storage

import "RestApi/internal/domain"

// TaskRepository — интерфейс для работы со списками
type TaskRepository interface {
	CreateTask(task domain.Task) (domain.Task, error)
	GetByIDTask(id string) (domain.Task, error)
	ListTasks(listID string, limit int, offset int) ([]domain.Task, int, error)
	UpdateTask(id string, text string, completed bool) (domain.Task, error)
	DeleteTask(id string) error
}
