package handlers

import (
	"net/http"
	"strconv"

	"github.com/gadzz/scopepay/internal/dto"
	"github.com/gadzz/scopepay/internal/httpx"
	"github.com/gadzz/scopepay/internal/models"
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

func (h *ProjectHandler) Get(c *gin.Context) {
	id64, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		httpx.Error(c, http.StatusBadRequest, httpx.ErrorBody{
			Code: httpx.CodeBadRequest, Message: "invalid project id",
		})
		return
	}
	p, err := h.svc.Get(c.Request.Context(), 1, uint(id64))

	if err != nil {
		httpx.WriteError(c, err)
		return
	}

	c.JSON(http.StatusOK, p)
}

func (h *ProjectHandler) List(c *gin.Context) {

	projects, err := h.svc.List(c.Request.Context(), 1)
	if err != nil {
		httpx.WriteError(c, err)
		return
	}

	if projects == nil {
		projects = []models.Project{}
	}
	c.JSON(200, projects)
}

func (h *ProjectHandler) Update(c *gin.Context) {

	id64, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		httpx.Error(c, http.StatusBadRequest, httpx.ErrorBody{
			Code: httpx.CodeBadRequest, Message: "invalid project id",
		})
		return
	}

	var req dto.UpdateProjectRequest
	if err := httpx.BindJSON(c, &req); err != nil {
		httpx.WriteError(c, err)
		return
	}

	p, err := h.svc.Update(c.Request.Context(), 1, uint(id64), req)

	if err != nil {
		httpx.WriteError(c, err)
		return
	}

	c.JSON(http.StatusOK, p)
}
