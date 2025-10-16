package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"RestApi/internal/service"
	"RestApi/internal/storage/mem"

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

type CreateListRequest struct {
	Title string `json:"title"`
}

type UpdateListRequest struct {
	Title string `json:"title"`
}

func (h *ListHandler) Create(w http.ResponseWriter, r *http.Request) {
	var request CreateListRequest
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		writeJSON(w, http.StatusBadRequest, ErrorResponse{
			Code:    "VALIDATION_FAILED",
			Message: "Invalid JSON format",
			Details: err.Error(),
		})
		return
	}

	list, err := h.service.Create(request.Title)
	if err != nil {
		if err == service.ErrValidation {
			writeJSON(w, http.StatusBadRequest, ErrorResponse{
				Code:    "VALIDATION_FAILED",
				Message: "title must be 1..100 chars",
				Details: err.Error(),
			})
			return
		}
		fmt.Printf("Error creating list: %v\n", err)

		writeJSON(w, http.StatusInternalServerError, ErrorResponse{
			Code:    "INTERNAL_ERROR",
			Message: "Internal server error",
			Details: err.Error(),
		})
		return
	}
	fmt.Printf("Создан список: ID=%s, Title=%q\n", list.ID, list.Title)

	writeJSON(w, http.StatusCreated, list)
}

func (h *ListHandler) GetByID(w http.ResponseWriter, r *http.Request) {

	params := mux.Vars(r)
	id := params["id"]

	list, err := h.service.GetByID(id)
	if err != nil {
		if err == mem.ErrNotFound {
			writeJSON(w, http.StatusNotFound, ErrorResponse{
				Code:    "NOT_FOUND",
				Message: "List not found",
				Details: err.Error(),
			})
			return
		}

		writeJSON(w, http.StatusInternalServerError, ErrorResponse{
			Code:    "INTERNAL_ERROR",
			Message: "Internal server error",
			Details: err.Error(),
		})
		return
	}

	writeJSON(w, http.StatusOK, list)
}

func (h *ListHandler) SearchByTitle(w http.ResponseWriter, r *http.Request) {

	query := r.URL.Query().Get("q")
	if query == "" {
		http.Error(w, "Query parameter 'q' is required", http.StatusBadRequest)
		return
	}

	lists, err := h.service.SearchByTitle(query)
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, ErrorResponse{
			Code:    "INTERNAL_ERROR",
			Message: "Internal server error",
			Details: err.Error(),
		})
		return
	}

	writeJSON(w, http.StatusOK, lists)
}

func (h *ListHandler) Update(w http.ResponseWriter, r *http.Request) {

	params := mux.Vars(r)
	id := params["id"]

	var request UpdateListRequest
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		writeJSON(w, http.StatusBadRequest, ErrorResponse{
			Code:    "VALIDATION_FAILED",
			Message: "Invalid JSON format",
			Details: err.Error(),
		})
		return
	}

	updatedList, err := h.service.Update(id, request.Title)
	if err != nil {
		if err == service.ErrValidation {
			writeJSON(w, http.StatusBadRequest, ErrorResponse{
				Code:    "VALIDATION_FAILED",
				Message: "title must be 1..100 chars",
				Details: err.Error(),
			})
			return
		}
		if err == mem.ErrNotFound {
			writeJSON(w, http.StatusNotFound, ErrorResponse{
				Code:    "NOT_FOUND",
				Message: "List not found",
				Details: err.Error(),
			})
			return
		}
		writeJSON(w, http.StatusInternalServerError, ErrorResponse{
			Code:    "INTERNAL_ERROR",
			Message: "Internal server error",
			Details: err.Error(),
		})
		return
	}
	writeJSON(w, http.StatusOK, updatedList)
}

func (h *ListHandler) Delete(w http.ResponseWriter, r *http.Request) {

	params := mux.Vars(r)
	id := params["id"]

	err := h.service.Delete(id)
	if err != nil {
		if err == mem.ErrNotFound {
			writeJSON(w, http.StatusNotFound, ErrorResponse{
				Code:    "NOT_FOUND",
				Message: "List not found",
				Details: err.Error(),
			})
			return
		}

		writeJSON(w, http.StatusInternalServerError, ErrorResponse{
			Code:    "INTERNAL_ERROR",
			Message: "Internal server error",
			Details: err.Error(),
		})
		return
	}

	writeJSON(w, http.StatusNoContent, nil)
}

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
		writeJSON(w, http.StatusInternalServerError, ErrorResponse{
			Code:    "INTERNAL_ERROR",
			Message: "Failed to paginate lists",
			Details: err.Error(),
		})
		return
	}
	w.Header().Set("X-Total-Count", strconv.Itoa(total))

	writeJSON(w, http.StatusOK, paginatedLists)
}

func writeJSON(w http.ResponseWriter, status int, data interface{}) {
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
