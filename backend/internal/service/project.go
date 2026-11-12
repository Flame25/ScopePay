package service

import (
	"context"
	"crypto/rand"
	"encoding/hex"

	"github.com/gadzz/scopepay/internal/dto"
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
