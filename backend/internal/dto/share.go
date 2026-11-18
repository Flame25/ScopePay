package dto

import "github.com/gadzz/scopepay/internal/models"

type ShareProjectResponse struct {
	Title        string               `json:"title"`
	ClientName   string               `json:"clientName"`
	ScopeSummary string               `json:"scopeSummary"`
	Status       models.ProjectStatus `json:"status"`
	Milestones   []models.Milestone   `json:"milestones"`
}

func NewShareProjectResponse(p *models.Project) ShareProjectResponse {
	ms := p.Milestones
	if ms == nil {
		ms = []models.Milestone{}
	}
	return ShareProjectResponse{
		Title:        p.Title,
		ClientName:   p.ClientName,
		ScopeSummary: p.ScopeSummary,
		Status:       p.Status,
		Milestones:   ms,
	}
}

type RequestChangesRequest struct {
	Note string `json:"note" binding:"max=2000"`
}
