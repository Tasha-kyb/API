package handlers

import "net/http"

// Health проверяет состояние сервиса
// @Summary Проверка здоровья
// @Description Проверяет, что сервис работает
// @Tags health
// @Produce json
// @Success 200 {object} map[string]string
// @Router /health [get]
func (h *ListHandler) Health(w http.ResponseWriter, r *http.Request) {
	WriteJSON(w, http.StatusOK, map[string]string{"status": "ok"})
}
