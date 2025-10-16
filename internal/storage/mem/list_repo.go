package mem
import (
	"errors"
	"time"
	"sync"

	"github.com/google/uuid"

	"RestApi/internal/domain"
)

type ListRepo struct {
	lists map[string]*domain.List
	mtx   sync.RWMutex
}

func NewListRepo() *ListRepo {
	return &ListRepo {
		lists: make(map[string]*domain.List),
	}
}

var (
	ErrListAlreadyExists = errors.New("list already exists")
	ErrNotFound = errors.New("NOT_FOUND")
)

func (l *ListRepo) Create(title string) (domain.List, error) {
	l.mtx.Lock()
	defer l.mtx.Unlock()

	id := uuid.NewString()

	if _, ok := l.lists[id]; ok {
		return domain.List{}, ErrListAlreadyExists
	}

	list := &domain.List {
		ID: 		id,
		Title: 		title,
		CreatedAt: 	time.Now().UTC(),
	}

	l.lists[id] = list
	return *list, nil
}

func (l *ListRepo) GetByID(id string) (domain.List, error) {
	l.mtx.RLock()
	defer l.mtx.RUnlock()

	list, ok := l.lists[id]
	if !ok {
		return domain.List{}, ErrNotFound
	}
	return *list, nil
}

func (l *ListRepo) Update(id string, title string) (domain.List, error) {
	l.mtx.Lock()
	defer l.mtx.Unlock()
	
	list, ok := l.lists[id] 
	if !ok {
		return domain.List{}, ErrNotFound
	}

	list.Title = title
	return *list, nil
}


func (l *ListRepo) Delete(id string) error {
	l.mtx.Lock()
	defer l.mtx.Unlock()

	if _, ok := l.lists[id]; !ok {
		return ErrNotFound
	}

	delete(l.lists, id)
	return nil
}

func (l *ListRepo) List(limit, offset int) ([]domain.List, int, error) {
	l.mtx.RLock()
	defer l.mtx.RUnlock()

	total := len(l.lists)
	if offset > total {
		return []domain.List{}, total, nil
	}

	all := make([]domain.List, 0, total)
	for _, list := range l.lists {
		all = append(all, *list)
	}

	start := offset
	end := total
	if limit > 0 && start+limit < end {
		end = start + limit
	}
	if start < 0 {
		start = 0
	}
	if start > end {
		start = end
	}
	return all[start:end], total, nil
}
