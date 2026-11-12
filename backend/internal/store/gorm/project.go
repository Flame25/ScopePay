package gorm

import (
	"context"
	"errors"

	"github.com/gadzz/scopepay/internal/httpx"
	"github.com/gadzz/scopepay/internal/models"
	"github.com/gadzz/scopepay/internal/store"
	pgorm "gorm.io/gorm"
)

type ProjectStore struct {
	db *pgorm.DB
}

func NewProjectStore(db *pgorm.DB) *ProjectStore {
	return &ProjectStore{db: db}
}

// Compile-time interface check:
// Ensures that *ProjectStore implements all methods required by store.ProjectStore.
// The compiler reports an error if any required method is missing or has an incompatible signature.
var _ store.ProjectStore = (*ProjectStore)(nil)

// ctx carries request-scoped information, such as cancellation and deadlines.
// Pass it through the application layers so database operations can respect
// the lifetime of the original HTTP request.
func (s *ProjectStore) Create(ctx context.Context, p *models.Project) error {
	return s.db.WithContext(ctx).Create(p).Error
}

func (s *ProjectStore) GetByID(ctx context.Context, id uint) (*models.Project, error) {
	var out models.Project
	err := s.db.WithContext(ctx).First(&out, id).Error
	if err != nil {
		if errors.Is(err, pgorm.ErrRecordNotFound) {
			return nil, httpx.ErrNotFound
		}
		return nil, err
	}
	return &out, nil
}

func (s *ProjectStore) ListByOwner(ctx context.Context, ownerID uint) ([]models.Project, error) {
	var out []models.Project
	err := s.db.WithContext(ctx).Where("freelancer_id = ?", ownerID).Order("created_at desc").Find(&out).Error
	return out, err
}
