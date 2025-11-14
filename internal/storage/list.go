package storage

import "RestApi/internal/domain"

// ListRepository — интерфейс для работы со списками
type ListRepository interface {
	Create(title string) (domain.List, error)
	GetByID(id string) (domain.List, error)
	SearchByTitle(title string) ([]domain.List, error)
	Update(id, title string) (domain.List, error)
	Delete(id string) error
	List(limit, offset int) ([]domain.List, int, error)
}
