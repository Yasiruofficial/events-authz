package http

import (
	"context"
	"errors"
	"github.com/spicedb/spicedb-go/internal/model"
	"github.com/spicedb/spicedb-go/internal/spicedb"

	"github.com/gin-gonic/gin"
)

type Handler struct {
	svc serviceChecker
}

type serviceChecker interface {
	Check(ctx context.Context, req model.CheckRequest) (model.CheckResponse, error)
}

func NewHandler(s serviceChecker) *Handler {
	return &Handler{svc: s}
}

func (h *Handler) Check(c *gin.Context) {
	var req model.CheckRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, gin.H{"error": "invalid request body"})
		return
	}

	allowed, err := h.svc.Check(c.Request.Context(), req)
	if err != nil {
		if errors.Is(err, spicedb.ErrInvalidCheckRequest) {
			c.JSON(400, gin.H{"error": err.Error()})
			return
		}

		c.JSON(502, gin.H{"error": "authorization backend error"})
		return
	}

	c.JSON(200, allowed)
}

func (h *Handler) Health(c *gin.Context) {
	c.String(200, "ok")
}
