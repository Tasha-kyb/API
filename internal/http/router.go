package http

import (
	"net/http"

	"github.com/gorilla/mux"

	"RestApi/internal/http/handlers"
)

type HTTPServer struct {
	router *mux.Router
}

func NewHTTPServer(httpHandler *handlers.ListHandler) *HTTPServer {
	router := mux.NewRouter()

	router.HandleFunc("/health", httpHandler.Health).Methods("GET")
	router.HandleFunc("/api/v1/lists", httpHandler.Create).Methods("POST")
	router.HandleFunc("/api/v1/lists", httpHandler.List).Methods("GET")
	router.HandleFunc("/api/v1/lists/search", httpHandler.SearchByTitle).Methods("GET")
	router.HandleFunc("/api/v1/lists/{id}", httpHandler.GetByID).Methods("GET")
	router.HandleFunc("/api/v1/lists/{id}", httpHandler.Update).Methods("PATCH")
	router.HandleFunc("/api/v1/lists/{id}", httpHandler.Delete).Methods("DELETE")

	return &HTTPServer{
		router: router,
	}
}

func (s *HTTPServer) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	s.router.ServeHTTP(w, r)
}
