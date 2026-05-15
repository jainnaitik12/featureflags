package controller

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus/promhttp"

	"featureflags/sdk-service/internal/service"
)

// SDKController exposes public read-only flag APIs and Prometheus metrics.
type SDKController struct {
	svc *service.EvalService
}

func NewSDKController(svc *service.EvalService) *SDKController {
	return &SDKController{svc: svc}
}

func (h *SDKController) RegisterRoutes(r *gin.Engine) {
	r.GET("/health", h.Health)
	r.GET("/metrics", gin.WrapH(promhttp.Handler()))
	r.GET("/flags", h.List)
	r.GET("/flags/:name", h.Get)
}

func (h *SDKController) Health(c *gin.Context) {
	c.String(http.StatusOK, "ok")
}

func (h *SDKController) List(c *gin.Context) {
	flags, err := h.svc.List(c.Request.Context())
	if err != nil {
		c.String(http.StatusInternalServerError, err.Error())
		return
	}
	c.JSON(http.StatusOK, flags)
}

func (h *SDKController) Get(c *gin.Context) {
	eval, err := h.svc.Get(c.Request.Context(), c.Param("name"))
	if errors.Is(err, service.ErrFlagNotFound) {
		c.String(http.StatusNotFound, "flag not found")
		return
	}
	if err != nil {
		c.String(http.StatusInternalServerError, err.Error())
		return
	}
	c.JSON(http.StatusOK, eval)
}
