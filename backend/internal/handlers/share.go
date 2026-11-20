package handlers

import (
	"net/http"

	"github.com/gadzz/scopepay/internal/dto"
	"github.com/gadzz/scopepay/internal/httpx"
	"github.com/gadzz/scopepay/internal/service"
	"github.com/gin-gonic/gin"
)

type ShareHandler struct {
	svc *service.ShareService
}

func NewShareHandler(svc *service.ShareService) *ShareHandler {
	return &ShareHandler{svc: svc}
}

func (h *ShareHandler) Get(c *gin.Context) {
	ownerID, ok := currentUserID(c)
	if !ok {
		return
	}

	p, role, err := h.svc.Get(c.Request.Context(), c.Param("token"), ownerID)
	if err != nil {
		httpx.WriteError(c, err)
		return
	}
	c.JSON(http.StatusOK, dto.NewShareProjectResponse(p, role))
}

func (h *ShareHandler) Approve(c *gin.Context) {
	ownerID, ok := currentUserID(c)
	if !ok {
		return
	}

	p, role, err := h.svc.Approve(c.Request.Context(), c.Param("token"), ownerID)
	if err != nil {
		httpx.WriteError(c, err)
		return
	}
	c.JSON(http.StatusOK, dto.NewShareProjectResponse(p, role))
}

func (h *ShareHandler) RequestChanges(c *gin.Context) {
	userID, ok := currentUserID(c)
	if !ok {
		return
	}

	var req dto.RequestChangesRequest
	if err := httpx.BindJSON(c, &req); err != nil {
		httpx.WriteError(c, err)
		return
	}

	p, role, err := h.svc.RequestChanges(c.Request.Context(), c.Param("token"), userID, req.Note)
	if err != nil {
		httpx.WriteError(c, err)
		return
	}
	c.JSON(http.StatusOK, dto.NewShareProjectResponse(p, role))
}
