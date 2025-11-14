// @title Tasks API
// @version 1.0.0
// @description REST API для управления списками и задачами (CRUD, PostgreSQL)

// @host localhost:8080
// @BasePath /

// @securityDefinitions.apikey ApiKeyAuth
// @in header
// @name Authorization

package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"RestApi/internal/config"
	"RestApi/internal/database"
	myhttp "RestApi/internal/http"
	"RestApi/internal/http/handlers"
	_ "RestApi/internal/http/handlers"
	"RestApi/internal/http/middleware"
	"RestApi/internal/service"
	"RestApi/internal/storage/postgres"
)

func main() {
	// Загружаем конфигурацию
	cfg := config.Load()

	// Создаем контекст для работы
	ctx := context.Background()

	// Подключаемся к PostgreSQL
	log.Println("Connecting to database...")
	pool, err := database.NewPool(ctx, cfg.DatabaseURL())
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer pool.Close()
	log.Println("Connected to database")

	// Создаем репозиторий PostgreSQL
	listRepo := postgres.NewListRepo(pool)
	taskRepo := postgres.NewTaskRepo(pool)

	// Создаем сервис
	listService := service.NewListService(listRepo)
	taskService := service.NewTaskService(taskRepo, listRepo)

	// Создаем HTTP-роутер
	listHandler := handlers.NewListHandler(listService)
	taskHandler := handlers.NewTaskHandler(taskService)

	httpServer := myhttp.NewHTTPServer(listHandler, taskHandler)

	// Создаем обработчик с middleware
	httpHandler := middleware.RequestID(httpServer)
	httpHandler = middleware.Logging(httpHandler)

	// Создаем HTTP-сервер
	server := &http.Server{
		Addr:         ":" + cfg.Port,
		Handler:      httpHandler,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	// Запускаем сервер в горутине
	go func() {
		log.Printf("Starting server on port %s...", cfg.Port)
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Server error: %v", err)
		}
	}()

	// Graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("Shutting down server...")
	shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := server.Shutdown(shutdownCtx); err != nil {
		log.Fatalf("Server forced to shutdown: %v", err)
	}

	log.Println("Server stopped")
}
