package handlers

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/gadzz/scopepay/internal/dto"
	"github.com/gadzz/scopepay/internal/httpx"
	"github.com/gadzz/scopepay/internal/models"
	"github.com/gadzz/scopepay/internal/service"
	"github.com/gin-gonic/gin"
)

type MilestoneHandler struct {
	svc *service.MilestoneService
}

func NewMilestoneHandler(svc *service.MilestoneService) *MilestoneHandler {
	return &MilestoneHandler{svc: svc}
}

func (h *MilestoneHandler) Create(c *gin.Context) {
	ownerID, ok := currentUserID(c)
	if !ok {
		return
	}

	var req dto.CreateMilestoneRequest
	err := httpx.BindJSON(c, &req)
	if err != nil {
		httpx.WriteError(c, err)
		return
	}

	projectID, ok := parseUintParam(c, "id")
	if !ok {
		return
	}

	p, err := h.svc.Create(c.Request.Context(), ownerID, projectID, req)
	if err != nil {
		httpx.WriteError(c, err)
		return
	}
	c.JSON(http.StatusCreated, p)

}

func (h *MilestoneHandler) List(c *gin.Context) {
	ownerID, ok := currentUserID(c)
	if !ok {
		return
	}

	projectID, ok := parseUintParam(c, "id")
	if !ok {
		return
	}

	projects, err := h.svc.List(c.Request.Context(), ownerID, projectID)
	if err != nil {
		httpx.WriteError(c, err)
		return
	}

	if projects == nil {
		projects = []models.Milestone{}
	}
	c.JSON(200, projects)
}

func (h *MilestoneHandler) Update(c *gin.Context) {
	ownerID, ok := currentUserID(c)
	if !ok {
		return
	}

	projectID, ok := parseUintParam(c, "id")
	if !ok {
		return
	}
	milestoneID, ok := parseUintParam(c, "mid")
	if !ok {
		return
	}

	var req dto.UpdateMilestoneRequest
	if err := httpx.BindJSON(c, &req); err != nil {
		httpx.WriteError(c, err)
		return
	}

	p, err := h.svc.Update(c.Request.Context(), ownerID, projectID, milestoneID, req)

	if err != nil {
		httpx.WriteError(c, err)
		return
	}

	c.JSON(http.StatusOK, p)
}

func (h *MilestoneHandler) Delete(c *gin.Context) {
	ownerID, ok := currentUserID(c)
	if !ok {
		return
	}

	projectID, ok := parseUintParam(c, "id")
	if !ok {
		return
	}
	milestoneID, ok := parseUintParam(c, "mid")
	if !ok {
		return
	}

	if err := h.svc.Delete(c.Request.Context(), ownerID, projectID, milestoneID); err != nil {
		httpx.WriteError(c, err)
		return
	}

	c.Status(http.StatusNoContent)

}

func parseUintParam(c *gin.Context, s string) (uint, bool) {

	data64, err := strconv.ParseUint(c.Param(s), 10, 64)
	if err != nil {
		httpx.Error(c, http.StatusBadRequest, httpx.ErrorBody{
			Code: httpx.CodeBadRequest, Message: fmt.Sprint("invalid project " + s),
		})
		return 0, false
	}

	res := uint(data64)

	return res, true
}
