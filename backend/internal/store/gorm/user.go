package gorm

import (
	"context"
	"errors"

	"github.com/gadzz/scopepay/internal/httpx"
	"github.com/gadzz/scopepay/internal/models"
	"github.com/gadzz/scopepay/internal/store"
	pgorm "gorm.io/gorm"
)

type UserStore struct {
	db *pgorm.DB
}

func NewUserStore(db *pgorm.DB) *UserStore {
	return &UserStore{db: db}
}

var _ store.UserStore = (*UserStore)(nil)

func (s *UserStore) Create(ctx context.Context, u *models.User) error {

	err := s.db.WithContext(ctx).Create(u).Error

	if errors.Is(err, pgorm.ErrDuplicatedKey) {
		return httpx.ErrConflict
	}

	return err
}

func (s *UserStore) GetByEmail(ctx context.Context, email string) (*models.User, error) {
	var out models.User

	err := s.db.WithContext(ctx).Where("email = ?", email).First(&out).Error

	if err != nil {
		if errors.Is(err, pgorm.ErrRecordNotFound) {
			return nil, httpx.ErrNotFound
		}
		return nil, err
	}

	return &out, nil
}

func (s *UserStore) GetByID(ctx context.Context, id uint) (*models.User, error) {
	var out models.User

	err := s.db.WithContext(ctx).First(&out, id).Error

	if err != nil {
		if errors.Is(err, pgorm.ErrRecordNotFound) {
			return nil, httpx.ErrNotFound
		}
		return nil, err
	}

	return &out, nil
}
