//go:build ignore
// +build ignore

package api_test

import (
	api "RestApi/internal/api"
	"testing"
)

func TestCodegenTypes(t *testing.T) {
	// Экземпляры сгенерированных типов
	task := api.Task{}
	createReq := api.CreateTaskRequest{Text: "Test"}
	completed := true
	updateReq := api.UpdateTaskRequest{Completed: &completed}

	t.Logf("Сгенерированные типы работают: %+v", task)
	t.Logf("CreateTaskRequest: %+v", createReq)
	t.Logf("UpdateTaskRequest: %+v", updateReq)
}
