package http

import (
	"net/http"

	"github.com/gorilla/mux"

	"RestApi/internal/http/handlers"

	_ "RestApi/docs"

	httpSwagger "github.com/swaggo/http-swagger"
)

type HTTPServer struct {
	router *mux.Router
}

func NewHTTPServer(httpHandler *handlers.ListHandler, taskHandlers *handlers.TaskHandler) *HTTPServer {
	router := mux.NewRouter()
	enableCORS(router)

	// Swagger UI - правильная настройка
	router.PathPrefix("/swagger/").Handler(httpSwagger.Handler(
		httpSwagger.URL("/swagger/doc.json"),
		httpSwagger.DeepLinking(true),
		httpSwagger.DocExpansion("none"),
		httpSwagger.DomID("swagger-ui"),
	))

	// Явно указываем endpoint для doc.json
	router.HandleFunc("/swagger/doc.json", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "./docs/swagger.json")
	})

	router.HandleFunc("/health", httpHandler.Health).Methods("GET")

	router.HandleFunc("/api/v1/lists", httpHandler.Create).Methods("POST")
	router.HandleFunc("/api/v1/lists", httpHandler.List).Methods("GET")
	router.HandleFunc("/api/v1/lists/search", httpHandler.SearchByTitle).Methods("GET")
	router.HandleFunc("/api/v1/lists/{id}", httpHandler.GetByID).Methods("GET")
	router.HandleFunc("/api/v1/lists/{id}", httpHandler.Update).Methods("PATCH")
	router.HandleFunc("/api/v1/lists/{id}", httpHandler.Delete).Methods("DELETE")

	router.HandleFunc("/api/v1/lists/{listID}/tasks", taskHandlers.CreateTask).Methods("POST")
	router.HandleFunc("/api/v1/lists/{listID}/tasks", taskHandlers.ListTasks).Methods("GET")
	router.HandleFunc("/api/v1/tasks/{taskID}", taskHandlers.GetTask).Methods("GET")
	router.HandleFunc("/api/v1/tasks/{taskID}", taskHandlers.UpdateTask).Methods("PATCH")
	router.HandleFunc("/api/v1/tasks/{taskID}", taskHandlers.DeleteTask).Methods("DELETE")

	return &HTTPServer{
		router: router,
	}
}

func (s *HTTPServer) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	s.router.ServeHTTP(w, r)
}

func enableCORS(router *mux.Router) {
	router.Methods("OPTIONS").HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
		w.WriteHeader(http.StatusOK)
	})

	router.Use(func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Access-Control-Allow-Origin", "*")
			next.ServeHTTP(w, r)
		})
	})
}
