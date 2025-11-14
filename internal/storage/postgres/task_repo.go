package postgres

import (
	"RestApi/internal/domain"
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type TaskRepo struct {
	pool *pgxpool.Pool
}

func NewTaskRepo(pool *pgxpool.Pool) *TaskRepo {
	return &TaskRepo{
		pool: pool,
	}
}

// Create создает новую задачу
func (r *TaskRepo) CreateTask(task domain.Task) (domain.Task, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Генерируем ID если не передан
	if task.ID == "" {
		task.ID = uuid.New().String()
	}
	// Устанавливаем временные метки если не установлены
	if task.CreatedAt.IsZero() {
		task.CreatedAt = time.Now()
	}
	if task.UpdatedAt.IsZero() {
		task.UpdatedAt = time.Now()
	}

	query := `
        INSERT INTO tasks (id, list_id, text, completed, created_at, updated_at)
        VALUES ($1, $2, $3, $4, $5, $6)
        RETURNING id, list_id, text, completed, created_at, updated_at
    `
	var createdTask domain.Task
	err := r.pool.QueryRow(ctx, query,
		task.ID,
		task.ListID,
		task.Text,
		task.Completed,
		task.CreatedAt,
		task.UpdatedAt,
	).Scan(
		&createdTask.ID,
		&createdTask.ListID,
		&createdTask.Text,
		&createdTask.Completed,
		&createdTask.CreatedAt,
		&createdTask.UpdatedAt,
	)
	if err != nil {
		return domain.Task{}, fmt.Errorf("create task: %w", err)
	}

	return createdTask, nil
}

// GetByID получает список по ID
func (r *TaskRepo) GetByIDTask(id string) (domain.Task, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	query := `
		SELECT id, list_id, text, completed, created_at, updated_at 
		FROM tasks 
		WHERE id = $1
	`

	var task domain.Task

	err := r.pool.QueryRow(ctx, query, id).Scan(
		&task.ID,
		&task.ListID,
		&task.Text,
		&task.Completed,
		&task.CreatedAt,
		&task.UpdatedAt,
	)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return domain.Task{}, ErrNotFound
		}
		return domain.Task{}, fmt.Errorf("get task by id: %w", err)
	}

	return task, nil
}

func (r *TaskRepo) ListTasks(listID string, limit int, offset int) ([]domain.Task, int, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Получаем общее количество задач в списке
	var total int
	countQuery := `SELECT COUNT(*) FROM tasks WHERE list_id = $1`
	err := r.pool.QueryRow(ctx, countQuery, listID).Scan(&total)
	if err != nil {
		return nil, 0, fmt.Errorf("count tasks: %w", err)
	}

	query := `
		SELECT id, list_id, text, completed, created_at, updated_at
		FROM tasks
		WHERE list_id = $1
		ORDER BY created_at DESC
		LIMIT $2 OFFSET $3
	`
	rows, err := r.pool.Query(ctx, query, listID, limit, offset)
	if err != nil {
		return nil, 0, fmt.Errorf("list tasks: %w", err)
	}
	defer rows.Close()

	tasks := make([]domain.Task, 0)
	for rows.Next() {
		var task domain.Task
		err := rows.Scan(
			&task.ID,
			&task.ListID,
			&task.Text,
			&task.Completed,
			&task.CreatedAt,
			&task.UpdatedAt,
		)
		if err != nil {
			return nil, 0, fmt.Errorf("scan task: %w", err)
		}
		tasks = append(tasks, task)
	}

	if err := rows.Err(); err != nil {
		return nil, 0, fmt.Errorf("rows error: %w", err)
	}

	return tasks, total, nil
}

// Update обновляет название задачи
func (r *TaskRepo) UpdateTask(id string, text string, completed bool) (domain.Task, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	query := `
        UPDATE tasks
		SET text = $2, completed = $3, updated_at = NOW()
		WHERE id = $1
		RETURNING id, list_id, text, completed, created_at, updated_at
    `

	var task domain.Task
	err := r.pool.QueryRow(ctx, query, id, text, completed).Scan(
		&task.ID,
		&task.ListID,
		&task.Text,
		&task.Completed,
		&task.CreatedAt,
		&task.UpdatedAt,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return domain.Task{}, ErrNotFound
		}
		return domain.Task{}, fmt.Errorf("update task: %w", err)
	}

	return task, nil
}

// Delete удаляет задачу
func (r *TaskRepo) DeleteTask(id string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	query := `DELETE FROM tasks WHERE id = $1`

	result, err := r.pool.Exec(ctx, query, id)
	if err != nil {
		return fmt.Errorf("delete task: %w", err)
	}

	if result.RowsAffected() == 0 {
		return ErrNotFound
	}

	return nil
}
