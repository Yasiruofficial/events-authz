package http

import (
	"events-authz/internal/model"
	"events-authz/internal/service"

	"github.com/gin-gonic/gin"
)

type Handler struct {
	svc *service.AuthzService
}

func NewHandler(s *service.AuthzService) *Handler {
	return &Handler{svc: s}
}

func (h *Handler) Check(c *gin.Context) {
	var req model.CheckRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.String(400, "invalid request")
		return
	}

	allowed, err := h.svc.Check(c.Request.Context(), req)
	if err != nil {
		c.String(500, "error")
		return
	}

	c.JSON(200, model.CheckResponse{Allowed: allowed})
}

func (h *Handler) Health(c *gin.Context) {
	c.String(200, "ok")
}
