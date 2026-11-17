package gorm

import (
	"context"
	"errors"

	"github.com/gadzz/scopepay/internal/httpx"
	"github.com/gadzz/scopepay/internal/models"
	"github.com/gadzz/scopepay/internal/store"
	pgorm "gorm.io/gorm"
)

type MilestoneStore struct {
	db *pgorm.DB
}

func NewMilestoneStore(db *pgorm.DB) *MilestoneStore {
	return &MilestoneStore{db: db}
}

var _ store.MilestoneStore = (*MilestoneStore)(nil)

func (s *MilestoneStore) Create(ctx context.Context, p *models.Milestone) error {
	return s.db.WithContext(ctx).Create(p).Error
}

func (s *MilestoneStore) GetByID(ctx context.Context, id uint) (*models.Milestone, error) {
	var out models.Milestone

	err := s.db.WithContext(ctx).First(&out, id).Error

	if err != nil {
		if errors.Is(err, pgorm.ErrRecordNotFound) {
			return nil, httpx.ErrNotFound
		}
		return nil, err
	}

	return &out, nil
}

func (s *MilestoneStore) Update(ctx context.Context, p *models.Milestone) error {
	return s.db.WithContext(ctx).Save(p).Error
}

func (s *MilestoneStore) Delete(ctx context.Context, mid uint) error {
	return s.db.WithContext(ctx).Delete(&models.Milestone{}, mid).Error
}

func (s *MilestoneStore) ListByProject(ctx context.Context, projectID uint) ([]models.Milestone, error) {
	var out []models.Milestone
	err := s.db.WithContext(ctx).Where("project_id = ?", projectID).Order("\"order\" asc").Find(&out).Error
	return out, err
}
