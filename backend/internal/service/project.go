package service

import (
	"context"
	"crypto/rand"
	"encoding/hex"

	"github.com/gadzz/scopepay/internal/dto"
	"github.com/gadzz/scopepay/internal/httpx"
	"github.com/gadzz/scopepay/internal/models"
	"github.com/gadzz/scopepay/internal/store"
)

type ProjectService struct {
	projects store.ProjectStore
}

func NewProjectService(projects store.ProjectStore) *ProjectService {
	return &ProjectService{projects: projects}
}

func newShareToken() (string, error) {
	var b [16]byte
	_, err := rand.Read(b[:])
	if err != nil {
		return "", err
	}
	return hex.EncodeToString(b[:]), nil
}

func (s *ProjectService) Create(ctx context.Context, ownerID uint, req dto.CreateProjectRequest) (*models.Project, error) {
	token, err := newShareToken()

	if err != nil {
		return nil, err
	}

	p := &models.Project{
		FreelancerID: ownerID,
		ClientName:   req.ClientName,
		ClientEmail:  req.ClientEmail,
		RawRequest:   req.RawRequest,
		Status:       models.ProjectDraft,
		ShareToken:   token,
	}

	err = s.projects.Create(ctx, p)
	if err != nil {
		return nil, err
	}

	return p, nil
}

func (s *ProjectService) Get(ctx context.Context, ownerID uint, id uint) (*models.Project, error) {
	p, err := s.projects.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if p.FreelancerID != ownerID {
		return nil, httpx.ErrForbidden
	}
	return p, nil
}

func (s *ProjectService) List(ctx context.Context, ownerID uint) ([]models.Project, error) {
	return s.projects.ListByOwner(ctx, ownerID)
}

func (s *ProjectService) Update(ctx context.Context, ownerID, id uint, req dto.UpdateProjectRequest) (*models.Project, error) {

	p, err := s.Get(ctx, ownerID, id)

	if err != nil {
		return nil, err
	}

	if p.Status != models.ProjectDraft && p.Status != models.ProjectSent {
		return nil, httpx.ErrConflict
	}

	if req.ClientName != nil {
		p.ClientName = *req.ClientName
	}

	if req.Title != nil {
		p.Title = *req.Title
	}

	if req.ClientEmail != nil {
		p.ClientEmail = *req.ClientEmail
	}

	if req.RawRequest != nil {
		p.RawRequest = *req.RawRequest
	}

	if req.ScopeSummary != nil {
		p.ScopeSummary = *req.ScopeSummary
	}

	if err := s.projects.Update(ctx, p); err != nil {
		return nil, err
	}

	return p, nil
}
