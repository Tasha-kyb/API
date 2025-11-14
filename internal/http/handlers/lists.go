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

type ListHandler struct {
	service *service.ListService
}

func NewListHandler(service *service.ListService) *ListHandler {
	return &ListHandler{
		service: service,
	}
}

// Create создает новый список
// @Summary Создать список
// @Description Создает новый список задач
// @Tags lists
// @Accept json
// @Produce json
// @Param input body domain.CreateListRequest true "Данные для создания списка"
// @Success 201 {object} domain.List
// @Failure 400 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /api/v1/lists [post]
func (h *ListHandler) Create(w http.ResponseWriter, r *http.Request) {
	var request domain.CreateListRequest
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		WriteJSON(w, http.StatusBadRequest, ErrorResponse{
			Code:    "VALIDATION_FAILED",
			Message: "Invalid JSON format",
			Details: err.Error(),
		})
		return
	}

	list, err := h.service.Create(request.Title)
	if err != nil {
		if err == service.ErrValidation {
			WriteJSON(w, http.StatusBadRequest, ErrorResponse{
				Code:    "VALIDATION_FAILED",
				Message: "title must be 1..100 chars",
				Details: err.Error(),
			})
			return
		}
		fmt.Printf("Error creating list: %v\n", err)

		WriteJSON(w, http.StatusInternalServerError, ErrorResponse{
			Code:    "INTERNAL_ERROR",
			Message: "Internal server error",
			Details: err.Error(),
		})
		return
	}
	fmt.Printf("Создан список: ID=%s, Title=%q\n", list.ID, list.Title)

	WriteJSON(w, http.StatusCreated, list)
}

// GetByID получает список по ID
// @Summary Получить список по ID
// @Description Возвращает список по его идентификатору
// @Tags lists
// @Accept json
// @Produce json
// @Param id path string true "ID списка"
// @Success 200 {object} domain.List
// @Failure 404 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /api/v1/lists/{id} [get]
func (h *ListHandler) GetByID(w http.ResponseWriter, r *http.Request) {

	params := mux.Vars(r)
	id := params["id"]

	list, err := h.service.GetByID(id)
	if err != nil {
		if err == postgres.ErrNotFound {
			WriteJSON(w, http.StatusNotFound, ErrorResponse{
				Code:    "NOT_FOUND",
				Message: "List not found",
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

	WriteJSON(w, http.StatusOK, list)
}

// SearchByTitle ищет списки по названию
// @Summary Поиск списков по названию
// @Description Возвращает списки, содержащие в названии заданную строку
// @Tags lists
// @Accept json
// @Produce json
// @Param q query string true "Поисковый запрос"
// @Success 200 {array} domain.List
// @Failure 400 {string} string "Неверный запрос"
// @Failure 500 {object} ErrorResponse
// @Router /api/v1/lists/search [get]
func (h *ListHandler) SearchByTitle(w http.ResponseWriter, r *http.Request) {

	query := r.URL.Query().Get("q")
	if query == "" {
		http.Error(w, "Query parameter 'q' is required", http.StatusBadRequest)
		return
	}

	lists, err := h.service.SearchByTitle(query)
	if err != nil {
		WriteJSON(w, http.StatusInternalServerError, ErrorResponse{
			Code:    "INTERNAL_ERROR",
			Message: "Internal server error",
			Details: err.Error(),
		})
		return
	}

	WriteJSON(w, http.StatusOK, lists)
}

// Update обновляет список
// @Summary Обновить список
// @Description Обновляет название списка
// @Tags lists
// @Accept json
// @Produce json
// @Param id path string true "ID списка"
// @Param input body domain.UpdateListRequest true "Новое название списка"
// @Success 200 {object} domain.List
// @Failure 400 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /api/v1/lists/{id} [patch]
func (h *ListHandler) Update(w http.ResponseWriter, r *http.Request) {

	params := mux.Vars(r)
	id := params["id"]

	var request domain.UpdateListRequest
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		WriteJSON(w, http.StatusBadRequest, ErrorResponse{
			Code:    "VALIDATION_FAILED",
			Message: "Invalid JSON format",
			Details: err.Error(),
		})
		return
	}

	updatedList, err := h.service.Update(id, request.Title)
	if err != nil {
		if err == service.ErrValidation {
			WriteJSON(w, http.StatusBadRequest, ErrorResponse{
				Code:    "VALIDATION_FAILED",
				Message: "title must be 1..100 chars",
				Details: err.Error(),
			})
			return
		}
		if err == postgres.ErrNotFound {
			WriteJSON(w, http.StatusNotFound, ErrorResponse{
				Code:    "NOT_FOUND",
				Message: "List not found",
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
	WriteJSON(w, http.StatusOK, updatedList)
}

// Delete удаляет список
// @Summary Удалить список
// @Description Удаляет список по его идентификатору
// @Tags lists
// @Accept json
// @Produce json
// @Param id path string true "ID списка"
// @Success 204 "Удалено"
// @Failure 404 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /api/v1/lists/{id} [delete]
func (h *ListHandler) Delete(w http.ResponseWriter, r *http.Request) {

	params := mux.Vars(r)
	id := params["id"]

	err := h.service.Delete(id)
	if err != nil {
		if err == postgres.ErrNotFound {
			WriteJSON(w, http.StatusNotFound, ErrorResponse{
				Code:    "NOT_FOUND",
				Message: "List not found",
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
		"message": "Успешное удаление",
	})
}

// List получает списки с пагинацией
// @Summary Получить списки
// @Description Возвращает список списков с пагинацией
// @Tags lists
// @Accept json
// @Produce json
// @Param limit query int false "Лимит" default(20)
// @Param offset query int false "Смещение" default(0)
// @Success 200 {array} domain.List
// @Header 200 {integer} X-Total-Count "Общее количество списков"
// @Failure 500 {object} ErrorResponse
// @Router /api/v1/lists [get]
func (h *ListHandler) List(w http.ResponseWriter, r *http.Request) {

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

	paginatedLists, total, err := h.service.List(limit, offset)

	if err != nil {
		WriteJSON(w, http.StatusInternalServerError, ErrorResponse{
			Code:    "INTERNAL_ERROR",
			Message: "Failed to paginate lists",
			Details: err.Error(),
		})
		return
	}
	w.Header().Set("X-Total-Count", strconv.Itoa(total))

	WriteJSON(w, http.StatusOK, paginatedLists)
}

func WriteJSON(w http.ResponseWriter, status int, data interface{}) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(status)

	if data != nil {
		if err := json.NewEncoder(w).Encode(data); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}

	}
}

type ErrorResponse struct {
	Code    string `json:"code"`
	Message string `json:"message"`
	Details string `json:"details"`
}
