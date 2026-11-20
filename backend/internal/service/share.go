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

const (
	RoleClient     = "client"
	RoleFreelancer = "freelancer"
)

func roleFor(p *models.Project, u *models.User) string {
	if p.FreelancerID == u.ID {
		return RoleFreelancer
	}
	return RoleClient
}

func NewShareService(projects store.ProjectStore, users store.UserStore) *ShareService {
	return &ShareService{projects: projects, users: users}
}

func (s *ShareService) Get(ctx context.Context, token string, userID uint) (*models.Project, string, error) {
	u, err := s.users.GetByID(ctx, userID)
	if err != nil {
		return nil, "", err
	}

	if token == "" {
		return nil, "", httpx.ErrNotFound
	}

	p, err := s.projects.GetByShareToken(ctx, token)
	if err != nil {
		return nil, "", err
	}

	if !strings.EqualFold(u.Email, p.ClientEmail) && p.FreelancerID != userID {
		return nil, "", httpx.ErrForbidden
	}

	return p, roleFor(p, u), nil
}

func (s *ShareService) Approve(ctx context.Context, token string, userID uint) (*models.Project, string, error) {
	p, role, err := s.Get(ctx, token, userID)
	if err != nil {
		return nil, "", err
	}

	if role != RoleClient {
		return nil, "", httpx.NewForbiddenError("only the client can approve this proposal")
	}

	switch p.Status {
	case models.ProjectSent:
		//ok
	case models.ProjectApproved:
		return p, role, nil
	default:
		return nil, "", httpx.NewStateError("only a sent proposal can be approved")
	}

	p.Status = models.ProjectApproved
	if err := s.projects.Update(ctx, p); err != nil {
		return nil, "", err
	}

	return p, role, nil
}

func (s *ShareService) RequestChanges(ctx context.Context, token string, userID uint, note string) (*models.Project, string, error) {
	p, role, err := s.Get(ctx, token, userID)
	if err != nil {
		return nil, "", err
	}

	switch p.Status {
	case models.ProjectSent, models.ProjectApproved, models.ProjectInProgress:
		// ok
	case models.ProjectRevision:
		return p, role, nil
	default:
		return nil, "", httpx.NewStateError("this proposal can no longer be reopened for changes")
	}

	p.Status = models.ProjectRevision
	p.EditRequestNote = note

	if err := s.projects.Update(ctx, p); err != nil {
		return nil, "", err
	}
	return p, role, nil
}
