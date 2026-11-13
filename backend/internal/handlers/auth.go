package handlers

import (
	"net/http"

	"github.com/gadzz/scopepay/internal/dto"
	"github.com/gadzz/scopepay/internal/httpx"
	"github.com/gadzz/scopepay/internal/service"
	"github.com/gin-gonic/gin"
)

type AuthHandler struct {
	svc *service.AuthService
}

func NewAuthHandler(svc *service.AuthService) *AuthHandler {
	return &AuthHandler{svc: svc}
}

func (h *AuthHandler) Register(c *gin.Context) {
	var req dto.RegisterRequest
	err := httpx.BindJSON(c, &req)

	if err != nil {
		httpx.WriteError(c, err)
		return
	}

	u, err := h.svc.Register(c.Request.Context(), req)

	if err != nil {
		httpx.WriteError(c, err)
		return
	}

	c.JSON(http.StatusCreated, u)
}

func (h *AuthHandler) Login(c *gin.Context) {
	var req dto.LoginRequest
	err := httpx.BindJSON(c, &req)

	if err != nil {
		httpx.WriteError(c, err)
		return
	}

	tok, u, err := h.svc.Login(c.Request.Context(), req)

	if err != nil {
		httpx.WriteError(c, err)
		return
	}

	c.JSON(http.StatusOK, dto.AuthResponse{Token: tok, User: *u})
}

func (h *AuthHandler) Me(c *gin.Context) {
	userId, ok := httpx.UserID(c)

	if !ok {
		httpx.Error(c, http.StatusUnauthorized, httpx.ErrorBody{Code: httpx.CodeUnauthorized, Message: "unauthenticated"})
		return
	}

	u, err := h.svc.Me(c.Request.Context(), userId)

	if err != nil {
		httpx.WriteError(c, err)
		return
	}
	c.JSON(http.StatusOK, u)
}
