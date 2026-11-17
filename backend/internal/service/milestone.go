package service

import (
	"context"

	"github.com/gadzz/scopepay/internal/dto"
	"github.com/gadzz/scopepay/internal/httpx"
	"github.com/gadzz/scopepay/internal/models"
	"github.com/gadzz/scopepay/internal/store"
)

type MilestoneService struct {
	projects   store.ProjectStore
	milestones store.MilestoneStore
}

func NewMilestoneService(milestones store.MilestoneStore, projects store.ProjectStore) *MilestoneService {
	return &MilestoneService{milestones: milestones, projects: projects}
}

func (s *MilestoneService) editableProject(ctx context.Context, ownerID, projectID uint) (*models.Project, error) {
	p, err := s.projects.GetByID(ctx, projectID)
	if err != nil {
		return nil, err
	}

	if p.FreelancerID != ownerID {
		return nil, httpx.ErrForbidden
	}

	if p.Status != models.ProjectDraft && p.Status != models.ProjectSent {
		return nil, httpx.ErrConflict
	}

	return p, nil
}

func (s *MilestoneService) Create(ctx context.Context, ownerID, projectID uint, req dto.CreateMilestoneRequest) (*models.Milestone, error) {
	if _, err := s.editableProject(ctx, ownerID, projectID); err != nil {
		return nil, err
	}

	existing, err := s.milestones.ListByProject(ctx, projectID)
	if err != nil {
		return nil, err
	}

	order := 0
	for _, m := range existing {
		if m.Order >= order {
			order = m.Order + 1
		}
	}

	currency := req.Currency
	if currency == "" {
		currency = "USD"
	}

	m := &models.Milestone{
		ProjectID:          projectID,
		Order:              order,
		Title:              req.Title,
		Deliverables:       req.Deliverables,
		AcceptanceCriteria: req.AcceptanceCriteria,
		AmountCents:        req.AmountCents,
		Currency:           currency,
		Status:             models.MilestonePending,
	}

	if err := s.milestones.Create(ctx, m); err != nil {
		return nil, err
	}
	return m, nil
}

func (s *MilestoneService) loadOwned(ctx context.Context, ownerID, projectID, milestoneID uint) (*models.Milestone, error) {
	if _, err := s.editableProject(ctx, ownerID, projectID); err != nil {
		return nil, err
	}
	m, err := s.milestones.GetByID(ctx, milestoneID)
	if err != nil {
		return nil, err
	}
	if m.ProjectID != projectID {
		return nil, httpx.ErrNotFound
	}
	return m, nil
}

func (s *MilestoneService) Update(ctx context.Context, ownerID, projectID, milestoneID uint, req dto.UpdateMilestoneRequest) (*models.Milestone, error) {
	m, err := s.loadOwned(ctx, ownerID, projectID, milestoneID)
	if err != nil {
		return nil, err
	}

	if m.Status != models.MilestonePending {
		return nil, httpx.ErrConflict
	}

	if req.Title != nil {
		m.Title = *req.Title
	}
	if req.Deliverables != nil {
		m.Deliverables = *req.Deliverables
	}
	if req.AcceptanceCriteria != nil {
		m.AcceptanceCriteria = *req.AcceptanceCriteria
	}
	if req.AmountCents != nil {
		m.AmountCents = *req.AmountCents
	}
	if req.Currency != nil {
		m.Currency = *req.Currency
	}
	if req.Order != nil {
		m.Order = *req.Order
	}

	if err := s.milestones.Update(ctx, m); err != nil {
		return nil, err
	}
	return m, nil
}

func (s *MilestoneService) Delete(ctx context.Context, ownerID, projectID, milestoneID uint) error {
	m, err := s.loadOwned(ctx, ownerID, projectID, milestoneID)
	if err != nil {
		return err
	}
	if m.Status != models.MilestonePending {
		return httpx.ErrConflict
	}
	return s.milestones.Delete(ctx, m.ID)
}

func (s *MilestoneService) List(ctx context.Context, ownerID, projectID uint) ([]models.Milestone, error) {
	p, err := s.projects.GetByID(ctx, projectID)
	if err != nil {
		return nil, err
	}
	if p.FreelancerID != ownerID {
		return nil, httpx.ErrForbidden
	}
	return s.milestones.ListByProject(ctx, projectID)
}
