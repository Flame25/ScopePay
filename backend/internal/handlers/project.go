package handlers

import (
	"net/http"

	"github.com/gadzz/scopepay/internal/dto"
	"github.com/gadzz/scopepay/internal/httpx"
	"github.com/gadzz/scopepay/internal/service"
	"github.com/gin-gonic/gin"
)

type ProjectHandler struct {
	svc *service.ProjectService
}

func NewProjectHandler(svc *service.ProjectService) *ProjectHandler {
	return &ProjectHandler{svc: svc}
}

func (h *ProjectHandler) Create(c *gin.Context) {
	var req dto.CreateProjectRequest
	err := httpx.BindJSON(c, &req)
	if err != nil {
		httpx.WriteError(c, err)
		return
	}

	p, err := h.svc.Create(c.Request.Context(), 1, req)
	if err != nil {
		httpx.WriteError(c, err)
		return
	}
	c.JSON(http.StatusCreated, p)
}
