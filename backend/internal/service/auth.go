package service

import (
	"context"
	"strings"

	"github.com/gadzz/scopepay/internal/auth"
	"github.com/gadzz/scopepay/internal/dto"
	"github.com/gadzz/scopepay/internal/httpx"
	"github.com/gadzz/scopepay/internal/models"
	"github.com/gadzz/scopepay/internal/store"
)

type AuthService struct {
	users  store.UserStore
	secret string
}

func NewAuthService(users store.UserStore, secret string) *AuthService {
	return &AuthService{users: users, secret: secret}
}

func (s *AuthService) Register(ctx context.Context, req dto.RegisterRequest) (*models.User, error) {

	email := strings.ToLower(strings.TrimSpace(req.Email))
	hash, err := auth.HashPassword(req.Password)

	if err != nil {
		return nil, err
	}

	u := &models.User{Email: email, Name: req.Name, PasswordHash: hash}
	err = s.users.Create(ctx, u)

	if err != nil {
		return nil, err
	}

	return u, nil

}

func (s *AuthService) Login(ctx context.Context, req dto.LoginRequest) (string, *models.User, error) {

	email := strings.ToLower(strings.TrimSpace(req.Email))
	u, err := s.users.GetByEmail(ctx, email)

	if err != nil {
		return "", nil, httpx.ErrUnauthorized
	}

	if !auth.CheckPassword(u.PasswordHash, req.Password) {
		return "", nil, httpx.ErrUnauthorized
	}

	tok, err := auth.IssueToken(s.secret, u.ID)

	if err != nil {
		return "", nil, err
	}
	return tok, u, nil
}

func (s *AuthService) Me(ctx context.Context, userID uint) (*models.User, error) {
	return s.users.GetByID(ctx, userID)
}
