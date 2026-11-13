package models

import "time"

type ProjectStatus string

const (
	ProjectDraft      ProjectStatus = "draft"
	ProjectSent       ProjectStatus = "sent"
	ProjectApproved   ProjectStatus = "approved"
	ProjectInProgress ProjectStatus = "in_progress"
	ProjectComplete   ProjectStatus = "complete"
)

type MilestoneStatus string

const (
	MilestonePending     MilestoneStatus = "pending"
	MilestoneAwaitingPay MilestoneStatus = "awaiting_payment"
	MilestonePaid        MilestoneStatus = "paid"
	MilestoneDelivered   MilestoneStatus = "delivered"
)

type PaymentStatus string

const (
	PaymentCreated  PaymentStatus = "created"
	PaymentApproved PaymentStatus = "approved"
	PaymentCaptured PaymentStatus = "captured"
	PaymentFailed   PaymentStatus = "failed"
)

type User struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	Email     string    `gorm:"uniqueIndex;not null" json:"email"`
	Name      string    `json:"name"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}

type Project struct {
	ID           uint          `gorm:"primaryKey" json:"id"`
	FreelancerID uint          `gorm:"index;not null" json:"freelancerId"`
	ClientName   string        `json:"clientName"`
	ClientEmail  string        `json:"clientEmail"`
	RawRequest   string        `gorm:"type:text" json:"rawRequest"`
	Title        string        `json:"title"`
	ScopeSummary string        `gorm:"type:text" json:"scopeSummary"`
	Status       ProjectStatus `gorm:"type:varchar(32);default:draft" json:"status"`
	ShareToken   string        `gorm:"uniqueIndex" json:"-"`
	Milestones   []Milestone   `json:"milestones,omitempty"`
	CreatedAt    time.Time     `json:"createdAt"`
	UpdatedAt    time.Time     `json:"updatedAt"`
}

type Milestone struct {
	ID                 uint            `gorm:"primaryKey" json:"id"`
	ProjectID          uint            `gorm:"index;not null" json:"projectId"`
	Order              int             `json:"order"`
	Title              string          `json:"title"`
	Deliverables       string          `gorm:"type:text" json:"deliverables"`
	AcceptanceCriteria string          `gorm:"type:text" json:"acceptanceCriteria"`
	AmountCents        int64           `json:"amountCents"`
	Currency           string          `gorm:"type:varchar(3);default:USD" json:"currency"`
	Status             MilestoneStatus `gorm:"type:varchar(32);default:pending" json:"status"`
	PaidAt             *time.Time      `json:"paidAt,omitempty"`
	CreatedAt          time.Time       `json:"createdAt"`
	UpdatedAt          time.Time       `json:"updatedAt"`
}

type Payment struct {
	ID              uint          `gorm:"primaryKey" json:"id"`
	MilestoneID     uint          `gorm:"index;not null" json:"milestoneId"`
	PayPalOrderID   string        `gorm:"uniqueIndex" json:"paypalOrderId"`
	PayPalCaptureID string        `gorm:"index" json:"paypalCaptureId"`
	AmountCents     int64         `json:"amountCents"`
	Status          PaymentStatus `gorm:"type:varchar(32);default:created" json:"status"`
	RawJSON         string        `gorm:"type:text" json:"-"`
	CreatedAt       time.Time     `json:"createdAt"`
	UpdatedAt       time.Time     `json:"updatedAt"`
}

type WebhookEvent struct {
	ID            uint       `gorm:"primaryKey" json:"id"`
	PayPalEventID string     `gorm:"uniqueIndex;not null" json:"paypalEventId"`
	Type          string     `json:"type"`
	Verified      bool       `json:"verified"`
	ProcessedAt   *time.Time `json:"processedAt,omitempty"`
	RawJSON       string     `gorm:"type:text" json:"-"`
	CreatedAt     time.Time  `json:"createdAt"`
}
