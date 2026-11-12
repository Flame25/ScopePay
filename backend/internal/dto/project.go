package dto

type CreateProjectRequest struct {
	ClientName  string `json:"clientName"  binding:"required,max=120"`
	ClientEmail string `json:"clientEmail" binding:"required,email,max=254"`
	RawRequest  string `json:"rawRequest"  binding:"required,min=10,max=5000"`
}
