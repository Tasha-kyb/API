//go:build integration
// +build integration

package postgres

import (
	"RestApi/internal/domain"
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/modules/postgres"
	"github.com/testcontainers/testcontainers-go/wait"
)

func setupTestDatabase(t *testing.T) *pgxpool.Pool {
	ctx := context.Background()

	container, err := postgres.Run(ctx,
		"postgres:15-alpine",
		postgres.WithDatabase("testdb"),
		postgres.WithUsername("user"),
		postgres.WithPassword("password"),
		testcontainers.WithWaitStrategy(
			wait.ForLog("database system is ready to accept connections").
				WithOccurrence(2).
				WithStartupTimeout(5*time.Second)),
	)
	require.NoError(t, err)

	connStr, err := container.ConnectionString(ctx)
	require.NoError(t, err)

	pool, err := pgxpool.New(ctx, connStr)
	require.NoError(t, err)

	// Create tables
	_, err = pool.Exec(ctx, `
		CREATE TABLE IF NOT EXISTS lists (
			id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
			title TEXT NOT NULL CHECK (length(title) BETWEEN 1 AND 100),
			created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
		)
	`)
	require.NoError(t, err)

	_, err = pool.Exec(ctx, `
		CREATE TABLE IF NOT EXISTS tasks (
			id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
			list_id UUID NOT NULL REFERENCES lists(id) ON DELETE CASCADE,
			text TEXT NOT NULL CHECK (length(text) BETWEEN 1 AND 500),
			completed BOOLEAN NOT NULL DEFAULT FALSE,
			created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
			updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
		)
	`)
	require.NoError(t, err)

	t.Cleanup(func() {
		pool.Close()
		container.Terminate(ctx)
	})

	return pool
}

func TestTaskRepo_Integration(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test")
	}

	pool := setupTestDatabase(t)
	repo := NewTaskRepo(pool)
	ctx := context.Background()

	// Create a test list first
	var listID string
	err := pool.QueryRow(ctx, "INSERT INTO lists (title) VALUES ($1) RETURNING id", "Test List").Scan(&listID)
	require.NoError(t, err)

	t.Run("Create and Get Task", func(t *testing.T) {
		task := domain.Task{
			ListID:    listID,
			Text:      "Integration test task",
			Completed: false,
		}

		created, err := repo.CreateTask(task)
		require.NoError(t, err)
		assert.NotEmpty(t, created.ID)
		assert.Equal(t, task.Text, created.Text)
		assert.Equal(t, task.ListID, created.ListID)
		assert.False(t, created.Completed)

		// Test GetByID
		fetched, err := repo.GetByIDTask(created.ID)
		require.NoError(t, err)
		assert.Equal(t, created, fetched)
	})

	t.Run("List Tasks with Pagination", func(t *testing.T) {
		// Create multiple tasks
		for i := 0; i < 5; i++ {
			_, err := repo.CreateTask(domain.Task{
				ListID: listID,
				Text:   fmt.Sprintf("Test task %d", i),
			})
			require.NoError(t, err)
		}

		tasks, total, err := repo.ListTasks(listID, 3, 0)
		require.NoError(t, err)
		assert.Len(t, tasks, 3)
		assert.GreaterOrEqual(t, total, 5)
	})

	t.Run("Update Task", func(t *testing.T) {
		task, _ := repo.CreateTask(domain.Task{
			ListID: listID,
			Text:   "To update",
		})

		updated, err := repo.UpdateTask(task.ID, "Updated text", true)
		require.NoError(t, err)
		assert.Equal(t, "Updated text", updated.Text)
		assert.True(t, updated.Completed)
	})

	t.Run("Delete Task", func(t *testing.T) {
		task, _ := repo.CreateTask(domain.Task{
			ListID: listID,
			Text:   "To delete",
		})

		err := repo.DeleteTask(task.ID)
		require.NoError(t, err)

		_, err = repo.GetByIDTask(task.ID)
		assert.Error(t, err)
	})
}
