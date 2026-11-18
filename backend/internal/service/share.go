package service

import (
	"context"
	"strings"

	"github.com/gadzz/scopepay/internal/httpx"
	"github.com/gadzz/scopepay/internal/models"
	"github.com/gadzz/scopepay/internal/store"
)

type ShareService struct {
	projects store.ProjectStore
	users    store.UserStore
}

func NewShareService(projects store.ProjectStore, users store.UserStore) *ShareService {
	return &ShareService{projects: projects, users: users}
}

func (s *ShareService) Get(ctx context.Context, token string, userID uint) (*models.Project, error) {
	if token == "" {
		return nil, httpx.ErrNotFound
	}

	p, err := s.projects.GetByShareToken(ctx, token)
	if err != nil {
		return nil, err
	}

	u, err := s.users.GetByID(ctx, userID)
	if err != nil {
		return nil, err
	}

	if !strings.EqualFold(u.Email, p.ClientEmail) && p.FreelancerID != userID {
		return nil, httpx.ErrForbidden
	}

	return p, nil

}

func (s *ShareService) Approve(ctx context.Context, token string, userID uint) (*models.Project, error) {
	p, err := s.Get(ctx, token, userID)
	if err != nil {
		return nil, err
	}

	switch p.Status {
	case models.ProjectSent:
		//ok
	case models.ProjectApproved:
		return p, nil
	default:
		return nil, httpx.ErrConflict
	}

	p.Status = models.ProjectApproved
	if err := s.projects.Update(ctx, p); err != nil {
		return nil, err
	}

	return p, nil
}

func (s *ShareService) RequestChanges(ctx context.Context, token string, userID uint, note string) (*models.Project, error) {
	p, err := s.Get(ctx, token, userID)
	if err != nil {
		return nil, err
	}

	switch p.Status {
	case models.ProjectApproved, models.ProjectInProgress:
		// ok
	case models.ProjectRevision:
		return p, nil
	default:
		return nil, httpx.NewStateError("only an approved proposal can be reopened for changes")
	}

	p.Status = models.ProjectRevision
	p.EditRequestNote = note

	if err := s.projects.Update(ctx, p); err != nil {
		return nil, err
	}
	return p, nil
}
