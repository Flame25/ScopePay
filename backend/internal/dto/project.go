package dto

type CreateProjectRequest struct {
	ClientName  string `json:"clientName"  binding:"required,max=120"`
	ClientEmail string `json:"clientEmail" binding:"required,email,max=254"`
	RawRequest  string `json:"rawRequest"  binding:"required,min=10,max=5000"`
}

type UpdateProjectRequest struct {
	ClientName   *string `json:"clientName"   binding:"omitempty,max=120"`
	ClientEmail  *string `json:"clientEmail"  binding:"omitempty,email,max=254"`
	RawRequest   *string `json:"rawRequest"   binding:"omitempty,min=10,max=5000"`
	Title        *string `json:"title"        binding:"omitempty,max=200"`
	ScopeSummary *string `json:"scopeSummary" binding:"omitempty,max=5000"`
}

type SendProjectResponse struct {
	ShareURL string `json:"shareUrl"`
}
