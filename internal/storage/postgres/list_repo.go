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

var ErrNotFound = errors.New("not found")

type ListRepo struct {
	pool    *pgxpool.Pool
	getByID string
}

func NewListRepo(pool *pgxpool.Pool) *ListRepo {
	return &ListRepo{
		pool:    pool,
		getByID: "SELECT id, title, created_at FROM lists WHERE id = $1",
	}
}

// Create создает новый список
func (r *ListRepo) Create(title string) (domain.List, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	id := uuid.New()
	query := `
        INSERT INTO lists (id, title)
        VALUES ($1, $2)
        RETURNING id, title, created_at
    `
	var list domain.List
	err := r.pool.QueryRow(ctx, query, id, title).Scan(
		&list.ID,
		&list.Title,
		&list.CreatedAt,
	)

	if err != nil {
		return domain.List{}, fmt.Errorf("create list: %w", err)
	}

	return list, nil
}

// GetByID получает список по ID
func (r *ListRepo) GetByID(id string) (domain.List, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var list domain.List

	err := r.pool.QueryRow(ctx, r.getByID, id).Scan(
		&list.ID,
		&list.Title,
		&list.CreatedAt,
	)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return domain.List{}, ErrNotFound
		}
		return domain.List{}, fmt.Errorf("get list by id: %w", err)
	}

	return list, nil
}

func (r *ListRepo) SearchByTitle(query string) ([]domain.List, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	searchQuery := `
        SELECT id, title, created_at 
        FROM lists 
        WHERE title ILIKE '%' || $1 || '%'
        ORDER BY created_at DESC
		`
	rows, err := r.pool.Query(ctx, searchQuery, query)
	if err != nil {
		return nil, fmt.Errorf("search lists by title: %w", err)
	}
	defer rows.Close()

	lists := make([]domain.List, 0)
	for rows.Next() {
		var list domain.List
		err := rows.Scan(&list.ID, &list.Title, &list.CreatedAt)
		if err != nil {
			return nil, fmt.Errorf("scan list: %w", err)
		}
		lists = append(lists, list)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("rows error: %w", err)
	}

	return lists, nil
}

// Update обновляет название списка
func (r *ListRepo) Update(id, title string) (domain.List, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	query := `
        UPDATE lists
        SET title = $2
        WHERE id = $1
        RETURNING id, title, created_at
    `

	var list domain.List
	err := r.pool.QueryRow(ctx, query, id, title).Scan(
		&list.ID,
		&list.Title,
		&list.CreatedAt,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return domain.List{}, ErrNotFound
		}
		return domain.List{}, fmt.Errorf("update list title: %w", err)
	}

	return list, nil
}

// Delete удаляет список
func (r *ListRepo) Delete(id string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	query := `DELETE FROM lists WHERE id = $1`

	result, err := r.pool.Exec(ctx, query, id)
	if err != nil {
		return fmt.Errorf("delete list: %w", err)
	}

	if result.RowsAffected() == 0 {
		return ErrNotFound
	}

	return nil
}

// List получает список с пагинацией
func (r *ListRepo) List(limit, offset int) ([]domain.List, int, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Получаем общее количество
	var total int
	countQuery := `SELECT COUNT(*) FROM lists`
	err := r.pool.QueryRow(ctx, countQuery).Scan(&total)
	if err != nil {
		return nil, 0, fmt.Errorf("count lists: %w", err)
	}

	// Получаем списки с пагинацией
	query := `
        SELECT id, title, created_at
        FROM lists
        ORDER BY created_at DESC
        LIMIT $1 OFFSET $2
    `

	rows, err := r.pool.Query(ctx, query, limit, offset)
	if err != nil {
		return nil, 0, fmt.Errorf("list lists: %w", err)
	}
	defer rows.Close()

	lists := make([]domain.List, 0)
	for rows.Next() {
		var list domain.List
		err := rows.Scan(&list.ID, &list.Title, &list.CreatedAt)
		if err != nil {
			return nil, 0, fmt.Errorf("scan list: %w", err)
		}
		lists = append(lists, list)
	}

	if err := rows.Err(); err != nil {
		return nil, 0, fmt.Errorf("rows error: %w", err)
	}

	return lists, total, nil
}

func (r *ListRepo) CreateWithItems(title string, items []string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Начинаем транзакцию
	tx, err := r.pool.Begin(ctx)
	if err != nil {
		return fmt.Errorf("begin transaction: %w", err)
	}
	defer tx.Rollback(ctx) // Откатим если что-то пойдет не так

	// Создаем список
	listID := uuid.New()
	_, err = tx.Exec(ctx,
		"INSERT INTO lists (id, title) VALUES ($1, $2)",
		listID, title,
	)
	if err != nil {
		return fmt.Errorf("create list: %w", err)
	}

	// Создаем элементы
	for _, item := range items {
		_, err = tx.Exec(ctx,
			"INSERT INTO items (list_id, text) VALUES ($1, $2)",
			listID, item,
		)
		if err != nil {
			return fmt.Errorf("create item: %w", err)
		}
	}

	// Коммитим транзакцию
	if err := tx.Commit(ctx); err != nil {
		return fmt.Errorf("commit transaction: %w", err)
	}

	return nil
}
