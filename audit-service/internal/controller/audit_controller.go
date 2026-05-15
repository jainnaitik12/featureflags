package controller

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"

	"featureflags/audit-service/internal/model"
	"featureflags/audit-service/internal/service"
)

// AuditController exposes audit HTTP routes.
type AuditController struct {
	svc *service.AuditService
}

func NewAuditController(svc *service.AuditService) *AuditController {
	return &AuditController{svc: svc}
}

func (h *AuditController) RegisterRoutes(r *gin.Engine) {
	r.GET("/health", h.Health)
	r.POST("/audit", h.Create)
	r.GET("/audit", h.List)
}

func (h *AuditController) Health(c *gin.Context) {
	c.String(http.StatusOK, "ok")
}

func (h *AuditController) Create(c *gin.Context) {
	var event model.AuditEvent
	if err := c.ShouldBindJSON(&event); err != nil {
		c.String(http.StatusBadRequest, "invalid json body")
		return
	}
	if err := h.svc.Record(c.Request.Context(), event); err != nil {
		if errors.Is(err, service.ErrValidation) {
			c.String(http.StatusBadRequest, err.Error())
			return
		}
		c.String(http.StatusInternalServerError, err.Error())
		return
	}
	c.Status(http.StatusCreated)
}

func (h *AuditController) List(c *gin.Context) {
	events, err := h.svc.Recent(c.Request.Context())
	if err != nil {
		c.String(http.StatusInternalServerError, err.Error())
		return
	}
	c.JSON(http.StatusOK, events)
}
