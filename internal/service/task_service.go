package service

import (
	"fmt"

	"RestApi/internal/domain"
	"RestApi/internal/storage"
	"RestApi/internal/storage/postgres"
)

type TaskService struct {
	repo     storage.TaskRepository
	listRepo storage.ListRepository
}

func NewTaskService(repo storage.TaskRepository, listRepo storage.ListRepository) *TaskService {
	return &TaskService{
		repo:     repo,
		listRepo: listRepo,
	}
}

func (l *TaskService) CreateTask(listID string, text string) (domain.Task, error) {
	if err := validateText(text); err != nil {
		return domain.Task{}, err
	}

	_, err := l.listRepo.GetByID(listID)
	if err != nil {
		if err == postgres.ErrNotFound {
			return domain.Task{}, fmt.Errorf("%w: list not found", ErrValidation)
		}
		return domain.Task{}, fmt.Errorf("failed to check list existence: %w", err)
	}

	task := domain.Task{
		ListID:    listID,
		Text:      text,
		Completed: false,
	}
	return l.repo.CreateTask(task)
}

func (l *TaskService) GetByIDTask(id string) (domain.Task, error) {
	return l.repo.GetByIDTask(id)
}

func (l *TaskService) ListTasks(listID string, limit int, offset int) ([]domain.Task, int, error) {
	return l.repo.ListTasks(listID, limit, offset)
}

func (l *TaskService) UpdateTask(id string, text *string, completed *bool) (domain.Task, error) {
	fmt.Printf("=== DEBUG UpdateTask Service ===\n")
	fmt.Printf("ID: %s\n", id)
	fmt.Printf("Text pointer: %v\n", text)
	if text != nil {
		fmt.Printf("Text value: '%s'\n", *text)
		fmt.Printf("Text length: %d\n", len(*text))
	}
	fmt.Printf("Completed pointer: %v\n", completed)
	fmt.Printf("==============================\n")
	// Получаем текущую задачу
	currentTask, err := l.repo.GetByIDTask(id)
	if err != nil {
		return domain.Task{}, err
	}

	// Обновляем текст только если передан
	newText := currentTask.Text
	if text != nil {
		if err := validateText(*text); err != nil {
			return domain.Task{}, err
		}
		newText = *text
	}

	// Обновляем статус только если передан
	newCompleted := currentTask.Completed
	if completed != nil {
		newCompleted = *completed
	}

	return l.repo.UpdateTask(id, newText, newCompleted)
}

func (l *TaskService) DeleteTask(id string) error {
	return l.repo.DeleteTask(id)
}

func validateText(text string) error {
	if len(text) == 0 || len(text) > 500 {
		return fmt.Errorf("%w: text must be 1..500 chars", ErrValidation)
	}
	return nil
}
