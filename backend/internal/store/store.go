package store

import (
	"context"

	"github.com/gadzz/scopepay/internal/models"
)

type ProjectStore interface {
	Create(ctx context.Context, p *models.Project) error
	ListByOwner(ctx context.Context, ownerID uint) ([]models.Project, error)
	GetByID(ctx context.Context, id uint) (*models.Project, error)
}
