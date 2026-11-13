package dto

import "github.com/gadzz/scopepay/internal/models"

type RegisterRequest struct {
	Email    string `json:"email"    binding:"required,email,max=254"`
	Name     string `json:"name"     binding:"required,max=120"`
	Password string `json:"password" binding:"required,min=8,max=72"`
}

type LoginRequest struct {
	Email    string `json:"email"    binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

type AuthResponse struct {
	Token string      `json:"token"`
	User  models.User `json:"user"`
}
