package service

import (
	"strings"
	"testing"

	"RestApi/internal/domain"
	"RestApi/internal/storage/postgres"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// Mock для TaskRepository
type MockTaskRepository struct {
	mock.Mock
}

func (m *MockTaskRepository) CreateTask(task domain.Task) (domain.Task, error) {
	args := m.Called(task)
	return args.Get(0).(domain.Task), args.Error(1)
}

func (m *MockTaskRepository) GetByIDTask(id string) (domain.Task, error) {
	args := m.Called(id)
	return args.Get(0).(domain.Task), args.Error(1)
}

func (m *MockTaskRepository) ListTasks(listID string, limit int, offset int) ([]domain.Task, int, error) {
	args := m.Called(listID, limit, offset)
	return args.Get(0).([]domain.Task), args.Int(1), args.Error(2)
}

func (m *MockTaskRepository) UpdateTask(id string, text string, completed bool) (domain.Task, error) {
	args := m.Called(id, text, completed)
	return args.Get(0).(domain.Task), args.Error(1)
}

func (m *MockTaskRepository) DeleteTask(id string) error {
	args := m.Called(id)
	return args.Error(0)
}

// Mock для ListRepository
type MockListRepository struct {
	mock.Mock
}

func (m *MockListRepository) Create(title string) (domain.List, error) {
	args := m.Called(title)
	return args.Get(0).(domain.List), args.Error(1)
}

func (m *MockListRepository) GetByID(id string) (domain.List, error) {
	args := m.Called(id)
	return args.Get(0).(domain.List), args.Error(1)
}

func (m *MockListRepository) SearchByTitle(query string) ([]domain.List, error) {
	args := m.Called(query)
	return args.Get(0).([]domain.List), args.Error(1)
}

func (m *MockListRepository) Update(id string, title string) (domain.List, error) {
	args := m.Called(id, title)
	return args.Get(0).(domain.List), args.Error(1)
}

func (m *MockListRepository) Delete(id string) error {
	args := m.Called(id)
	return args.Error(0)
}

func (m *MockListRepository) List(limit, offset int) ([]domain.List, int, error) {
	args := m.Called(limit, offset)
	return args.Get(0).([]domain.List), args.Int(1), args.Error(2)
}

func TestTaskService_CreateTask_Success(t *testing.T) {
	// Создаем моки
	taskRepo := new(MockTaskRepository)
	listRepo := new(MockListRepository)
	service := NewTaskService(taskRepo, listRepo)

	// Настраиваем ожидания:
	// - При проверке списка вернуть успех
	// - При создании задачи вернуть созданную задачу
	listRepo.On("GetByID", "list-123").Return(domain.List{ID: "list-123", Title: "Test List"}, nil)
	taskRepo.On("CreateTask", mock.AnythingOfType("domain.Task")).
		Return(domain.Task{
			ID:        "task-123",
			ListID:    "list-123",
			Text:      "Test task",
			Completed: false,
		}, nil)

	// Вызываем метод
	result, err := service.CreateTask("list-123", "Test task")

	// Проверяем результат
	assert.NoError(t, err)
	assert.Equal(t, "list-123", result.ListID)
	assert.Equal(t, "Test task", result.Text)
	assert.False(t, result.Completed)

	// Проверяем что все ожидаемые вызовы произошли
	taskRepo.AssertExpectations(t)
	listRepo.AssertExpectations(t)
}

func TestTaskService_CreateTask_EmptyText(t *testing.T) {
	taskRepo := new(MockTaskRepository)
	listRepo := new(MockListRepository)
	service := NewTaskService(taskRepo, listRepo)

	// Не настраиваем вызовы к репозиториям - их не должно быть при ошибке валидации

	// Вызываем метод с пустым текстом
	_, err := service.CreateTask("list-123", "")

	// Проверяем что получили ошибку валидации
	assert.Error(t, err)
	assert.ErrorIs(t, err, ErrValidation)
}

func TestTaskService_CreateTask_ListNotFound(t *testing.T) {
	taskRepo := new(MockTaskRepository)
	listRepo := new(MockListRepository)
	service := NewTaskService(taskRepo, listRepo)

	// Настраиваем что список не найден
	listRepo.On("GetByID", "non-existent-list").Return(domain.List{}, postgres.ErrNotFound)

	// Вызываем метод
	_, err := service.CreateTask("non-existent-list", "Test task")

	// Проверяем что получили ошибку
	assert.Error(t, err)
	assert.ErrorIs(t, err, ErrValidation)

	// Проверяем что создание задачи НЕ вызывалось
	taskRepo.AssertNotCalled(t, "CreateTask")
	listRepo.AssertExpectations(t)
}

func TestTaskService_UpdateTask_Success(t *testing.T) {
	taskRepo := new(MockTaskRepository)
	listRepo := new(MockListRepository)
	service := NewTaskService(taskRepo, listRepo)

	// Настраиваем мок для получения текущей задачи
	taskRepo.On("GetByIDTask", "task-123").
		Return(domain.Task{
			ID:        "task-123",
			ListID:    "list-123",
			Text:      "Original text",
			Completed: false,
		}, nil)

	// Настраиваем успешное обновление
	taskRepo.On("UpdateTask", "task-123", "Updated text", true).
		Return(domain.Task{
			ID:        "task-123",
			ListID:    "list-123",
			Text:      "Updated text",
			Completed: true,
		}, nil)

	text := "Updated text"
	completed := true
	result, err := service.UpdateTask("task-123", &text, &completed)

	assert.NoError(t, err)
	assert.Equal(t, "Updated text", result.Text)
	assert.True(t, result.Completed)
	taskRepo.AssertExpectations(t)
}

func TestTaskService_DeleteTask_Success(t *testing.T) {
	taskRepo := new(MockTaskRepository)
	listRepo := new(MockListRepository)
	service := NewTaskService(taskRepo, listRepo)

	// Настраиваем успешное удаление
	taskRepo.On("DeleteTask", "task-123").Return(nil)

	err := service.DeleteTask("task-123")

	assert.NoError(t, err)
	taskRepo.AssertExpectations(t)
}

func TestTaskService_EdgeCases(t *testing.T) {
	t.Run("text exactly 500 characters", func(t *testing.T) {
		taskRepo := new(MockTaskRepository)
		listRepo := new(MockListRepository)
		service := NewTaskService(taskRepo, listRepo)

		// Текст ровно 500 символов - должен работать
		maxText := strings.Repeat("a", 500)

		listRepo.On("GetByID", "list-123").Return(domain.List{ID: "list-123"}, nil)
		taskRepo.On("CreateTask", mock.Anything).Return(domain.Task{ID: "task-123"}, nil)

		_, err := service.CreateTask("list-123", maxText)
		assert.NoError(t, err)
	})

	t.Run("update with only completed flag", func(t *testing.T) {
		taskRepo := new(MockTaskRepository)
		listRepo := new(MockListRepository)
		service := NewTaskService(taskRepo, listRepo)

		// Настраиваем мок для получения текущей задачи
		taskRepo.On("GetByIDTask", "task-123").
			Return(domain.Task{
				ID:        "task-123",
				ListID:    "list-123",
				Text:      "Original text",
				Completed: false,
			}, nil)

		// Обновляем только completed, text остается прежним
		taskRepo.On("UpdateTask", "task-123", "Original text", true).
			Return(domain.Task{
				ID:        "task-123",
				Text:      "Original text",
				Completed: true,
			}, nil)

		completed := true
		_, err := service.UpdateTask("task-123", nil, &completed)
		assert.NoError(t, err)
	})
}
