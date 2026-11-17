package dto

type CreateMilestoneRequest struct {
	Title              string `json:"title"              binding:"required,max=200"`
	Deliverables       string `json:"deliverables"       binding:"max=2000"`
	AcceptanceCriteria string `json:"acceptanceCriteria" binding:"max=2000"`
	AmountCents        int64  `json:"amountCents"        binding:"required,gt=0"`
	Currency           string `json:"currency"           binding:"omitempty,oneof=USD EUR GBP"`
}

type UpdateMilestoneRequest struct {
	Title              *string `json:"title"              binding:"omitempty,max=200"`
	Deliverables       *string `json:"deliverables"       binding:"omitempty,max=2000"`
	AcceptanceCriteria *string `json:"acceptanceCriteria" binding:"omitempty,max=2000"`
	AmountCents        *int64  `json:"amountCents"        binding:"omitempty,gt=0"`
	Currency           *string `json:"currency"           binding:"omitempty,oneof=USD EUR GBP"`
	Order              *int    `json:"order"              binding:"omitempty,gte=0"`
}
