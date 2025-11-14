package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"RestApi/internal/domain"
	"RestApi/internal/service"
	"RestApi/internal/storage/postgres"

	"github.com/gorilla/mux"
)

type TaskHandler struct {
	service *service.TaskService
}

func NewTaskHandler(service *service.TaskService) *TaskHandler {
	return &TaskHandler{
		service: service,
	}
}

// CreateTask создает новую задачу
// @Summary Создать задачу
// @Description Создает новую задачу в указанном списке
// @Tags tasks
// @Accept json
// @Produce json
// @Param listID path string true "ID списка"
// @Param input body domain.CreateTaskRequest true "Данные для создания задачи"
// @Success 201 {object} domain.Task
// @Failure 400 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /api/v1/lists/{listID}/tasks [post]
func (h *TaskHandler) CreateTask(w http.ResponseWriter, r *http.Request) {

	params := mux.Vars(r)
	listID := params["listID"]

	var request domain.CreateTaskRequest

	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		WriteJSON(w, http.StatusBadRequest, ErrorResponse{
			Code:    "VALIDATION_FAILED",
			Message: "Invalid JSON format",
			Details: err.Error(),
		})
		return
	}

	task, err := h.service.CreateTask(listID, request.Text)
	if err != nil {
		if err == service.ErrValidation {
			WriteJSON(w, http.StatusBadRequest, ErrorResponse{
				Code:    "VALIDATION_FAILED",
				Message: "text must be 1..500 chars",
				Details: err.Error(),
			})
			return
		}
		fmt.Printf("Error creating task: %v\n", err)

		WriteJSON(w, http.StatusInternalServerError, ErrorResponse{
			Code:    "INTERNAL_ERROR",
			Message: "Internal server error",
			Details: err.Error(),
		})
		return
	}
	fmt.Printf("Создана задача: ID=%s, ListID=%s, Text=%q\n", task.ID, task.ListID, task.Text)

	WriteJSON(w, http.StatusCreated, task)
}

// GetTask получает задачу по ID
// @Summary Получить задачу по ID
// @Description Возвращает задачу по ее идентификатору
// @Tags tasks
// @Accept json
// @Produce json
// @Param taskID path string true "ID задачи"
// @Success 200 {object} domain.Task
// @Failure 404 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /api/v1/tasks/{taskID} [get]
func (h *TaskHandler) GetTask(w http.ResponseWriter, r *http.Request) {

	params := mux.Vars(r)
	id := params["taskID"]

	task, err := h.service.GetByIDTask(id)
	if err != nil {
		if err == postgres.ErrNotFound {
			WriteJSON(w, http.StatusNotFound, ErrorResponse{
				Code:    "NOT_FOUND",
				Message: "Task not found",
				Details: err.Error(),
			})
			return
		}

		WriteJSON(w, http.StatusInternalServerError, ErrorResponse{
			Code:    "INTERNAL_ERROR",
			Message: "Internal server error",
			Details: err.Error(),
		})
		return
	}

	WriteJSON(w, http.StatusOK, task)
}

// ListTasks получает задачи списка
// @Summary Получить задачи списка
// @Description Возвращает задачи указанного списка с пагинацией
// @Tags tasks
// @Accept json
// @Produce json
// @Param listID path string true "ID списка"
// @Param limit query int false "Лимит" default(20)
// @Param offset query int false "Смещение" default(0)
// @Success 200 {array} domain.Task
// @Header 200 {integer} X-Total-Count "Общее количество задач"
// @Failure 404 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /api/v1/lists/{listID}/tasks [get]
func (h *TaskHandler) ListTasks(w http.ResponseWriter, r *http.Request) {

	params := mux.Vars(r)
	listID := params["listID"]

	limitStr := r.URL.Query().Get("limit")
	offsetStr := r.URL.Query().Get("offset")

	limit := 20
	if limitStr != "" {
		if l, err := strconv.Atoi(limitStr); err == nil && l >= 0 {
			limit = l
		}
	}

	if limit > 100 {
		limit = 100
	}

	offset := 0
	if offsetStr != "" {
		if o, err := strconv.Atoi(offsetStr); err == nil && o >= 0 {
			offset = o
		}
	}

	tasks, total, err := h.service.ListTasks(listID, limit, offset)
	if err != nil {
		WriteJSON(w, http.StatusInternalServerError, ErrorResponse{
			Code:    "INTERNAL_ERROR",
			Message: "Failed to get tasks",
			Details: err.Error(),
		})
		return
	}

	w.Header().Set("X-Total-Count", strconv.Itoa(total))
	WriteJSON(w, http.StatusOK, tasks)
}

// Update обновляет задачу
// @Summary Обновить задачу
// @Description Обновляет описание и/или статус выполнения задачи
// @Tags tasks
// @Accept json
// @Produce json
// @Param taskID path string true "ID списка"
// @Param input body domain.UpdateTaskRequest true "Данные для обновления задачи"
// @Success 200 {object} domain.Task
// @Failure 400 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /api/v1/tasks/{taskID} [patch]
func (h *TaskHandler) UpdateTask(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	taskID := params["taskID"]

	var request domain.UpdateTaskRequest
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		WriteJSON(w, http.StatusBadRequest, ErrorResponse{
			Code:    "VALIDATION_FAILED",
			Message: "Invalid JSON format",
			Details: err.Error(),
		})
		return
	}

	fmt.Printf("=== DEBUG UpdateTask Handler ===\n")
	fmt.Printf("TaskID: %s\n", taskID)
	fmt.Printf("Request Text: %v\n", request.Text)
	if request.Text != nil {
		fmt.Printf("Request Text Value: '%s'\n", *request.Text)
		fmt.Printf("Request Text Length: %d\n", len(*request.Text))
	}
	fmt.Printf("Request Completed: %v\n", request.Completed)
	if request.Completed != nil {
		fmt.Printf("Request Completed Value: %v\n", *request.Completed)
	}
	fmt.Printf("===============================\n")

	if request.Text == nil && request.Completed == nil {
		WriteJSON(w, http.StatusBadRequest, ErrorResponse{
			Code:    "VALIDATION_FAILED",
			Message: "At least one field (text or completed) must be provided",
			Details: "No fields to update",
		})
		return
	}

	updatedTask, err := h.service.UpdateTask(taskID, request.Text, request.Completed)
	if err != nil {
		if err == service.ErrValidation {
			WriteJSON(w, http.StatusBadRequest, ErrorResponse{
				Code:    "VALIDATION_FAILED",
				Message: "text must be 1..500 chars",
				Details: err.Error(),
			})
			return
		}
		if err == postgres.ErrNotFound {
			WriteJSON(w, http.StatusNotFound, ErrorResponse{
				Code:    "NOT_FOUND",
				Message: "Task not found",
				Details: err.Error(),
			})
			return
		}
		WriteJSON(w, http.StatusInternalServerError, ErrorResponse{
			Code:    "INTERNAL_ERROR",
			Message: "Internal server error",
			Details: err.Error(),
		})
		return
	}
	WriteJSON(w, http.StatusOK, updatedTask)
}

// Delete удаляет задачу
// @Summary Удалить задачу
// @Description Удаляет задачу по ее идентификатору
// @Tags tasks
// @Accept json
// @Produce json
// @Param taskID path string true "ID задачи"
// @Success 204 "Удалено"
// @Failure 404 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /api/v1/tasks/{taskID} [delete]
func (h *TaskHandler) DeleteTask(w http.ResponseWriter, r *http.Request) {

	params := mux.Vars(r)
	taskID := params["taskID"]

	err := h.service.DeleteTask(taskID)
	if err != nil {
		if err == postgres.ErrNotFound {
			WriteJSON(w, http.StatusNotFound, ErrorResponse{
				Code:    "NOT_FOUND",
				Message: "Task not found",
				Details: err.Error(),
			})
			return
		}

		WriteJSON(w, http.StatusInternalServerError, ErrorResponse{
			Code:    "INTERNAL_ERROR",
			Message: "Internal server error",
			Details: err.Error(),
		})
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{
		"message": "Задача успешно удалена",
	})
}
