package service

import (
	"errors"
	"fmt"

	"RestApi/internal/domain"
	"RestApi/internal/storage"
)

var ErrValidation = errors.New("VALIDATION_FAILED")

type ListService struct {
	repo storage.ListRepository
}

func NewListService(repo storage.ListRepository) *ListService {
	return &ListService{repo: repo}
}

func (l *ListService) Create(title string) (domain.List, error) {
	if err := validateTitle(title); err != nil {
		return domain.List{}, err
	}
	return l.repo.Create(title)
}

func (l *ListService) GetByID(id string) (domain.List, error) {
	return l.repo.GetByID(id)
}

func (l *ListService) SearchByTitle(query string) ([]domain.List, error) {
	return l.repo.SearchByTitle(query)
}

func (l *ListService) Update(id string, title string) (domain.List, error) {

	if err := validateTitle(title); err != nil {
		return domain.List{}, err
	}
	return l.repo.Update(id, title)
}

func (l *ListService) Delete(id string) error {
	return l.repo.Delete(id)
}

func (l *ListService) List(limit, offset int) ([]domain.List, int, error) {
	return l.repo.List(limit, offset)
}

func validateTitle(title string) error {
	if len(title) == 0 || len(title) > 100 {
		return fmt.Errorf("%w: title must be 1..100 chars", ErrValidation)
	}
	return nil
}
