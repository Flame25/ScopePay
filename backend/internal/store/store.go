package store

import (
	"context"

	"github.com/gadzz/scopepay/internal/models"
)

type ProjectStore interface {
	Create(ctx context.Context, p *models.Project) error
	ListByOwner(ctx context.Context, ownerID uint) ([]models.Project, error)
	GetByID(ctx context.Context, id uint) (*models.Project, error)
	Update(ctx context.Context, p *models.Project) error
}

type UserStore interface {
	Create(ctx context.Context, u *models.User) error
	GetByEmail(ctx context.Context, email string) (*models.User, error)
	GetByID(ctx context.Context, id uint) (*models.User, error)
}
