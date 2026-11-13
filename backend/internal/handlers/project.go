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
	svc           *service.ProjectService
	publicBaseURL string
}

func NewProjectHandler(svc *service.ProjectService, publicBaseURL string) *ProjectHandler {
	return &ProjectHandler{svc: svc, publicBaseURL: publicBaseURL}
}

func (h *ProjectHandler) Create(c *gin.Context) {
	ownerID, ok := currentUserID(c)
	if !ok {
		return
	}

	var req dto.CreateProjectRequest
	err := httpx.BindJSON(c, &req)
	if err != nil {
		httpx.WriteError(c, err)
		return
	}

	p, err := h.svc.Create(c.Request.Context(), ownerID, req)
	if err != nil {
		httpx.WriteError(c, err)
		return
	}
	c.JSON(http.StatusCreated, p)
}

func (h *ProjectHandler) Get(c *gin.Context) {
	ownerID, ok := currentUserID(c)
	if !ok {
		return
	}

	id64, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		httpx.Error(c, http.StatusBadRequest, httpx.ErrorBody{
			Code: httpx.CodeBadRequest, Message: "invalid project id",
		})
		return
	}
	p, err := h.svc.Get(c.Request.Context(), ownerID, uint(id64))

	if err != nil {
		httpx.WriteError(c, err)
		return
	}

	c.JSON(http.StatusOK, p)
}

func (h *ProjectHandler) List(c *gin.Context) {
	ownerID, ok := currentUserID(c)
	if !ok {
		return
	}

	projects, err := h.svc.List(c.Request.Context(), ownerID)
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
	ownerID, ok := currentUserID(c)
	if !ok {
		return
	}

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

	p, err := h.svc.Update(c.Request.Context(), ownerID, uint(id64), req)

	if err != nil {
		httpx.WriteError(c, err)
		return
	}

	c.JSON(http.StatusOK, p)
}

func (h *ProjectHandler) Send(c *gin.Context) {
	ownerID, ok := currentUserID(c)
	if !ok {
		return
	}

	id64, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		httpx.Error(c, http.StatusBadRequest, httpx.ErrorBody{
			Code: httpx.CodeBadRequest, Message: "invalid project id",
		})
		return
	}
	p, err := h.svc.Send(c.Request.Context(), ownerID, uint(id64))

	if err != nil {
		httpx.WriteError(c, err)
		return
	}

	c.JSON(http.StatusOK, dto.SendProjectResponse{ShareURL: h.publicBaseURL + "/share/" + p.ShareToken})
}

func currentUserID(c *gin.Context) (uint, bool) {
	id, ok := httpx.UserID(c)
	if !ok {
		httpx.Error(c, http.StatusUnauthorized, httpx.ErrorBody{
			Code: httpx.CodeUnauthorized, Message: "unauthenticated",
		})
	}
	return id, ok
}
