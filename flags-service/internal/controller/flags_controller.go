package controller

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"

	"featureflags/flags-service/internal/model"
	"featureflags/flags-service/internal/service"
)

// FlagsController wires HTTP routes to FlagService.
type FlagsController struct {
	svc *service.FlagService
}

func NewFlagsController(svc *service.FlagService) *FlagsController {
	return &FlagsController{svc: svc}
}

func (h *FlagsController) RegisterRoutes(r *gin.Engine) {
	r.GET("/health", h.Health)
	r.GET("/flags", h.List)
	r.POST("/flags", h.Create)
	r.GET("/flags/:name", h.Get)
	r.PATCH("/flags/:name/toggle", h.Toggle)
	r.DELETE("/flags/:name", h.Delete)
}

func (h *FlagsController) Health(c *gin.Context) {
	c.String(http.StatusOK, "ok")
}

func (h *FlagsController) List(c *gin.Context) {
	flags, err := h.svc.List(c.Request.Context())
	if err != nil {
		c.String(http.StatusInternalServerError, err.Error())
		return
	}
	c.JSON(http.StatusOK, flags)
}

func (h *FlagsController) Get(c *gin.Context) {
	flag, err := h.svc.Get(c.Request.Context(), c.Param("name"))
	if errors.Is(err, service.ErrFlagNotFound) {
		c.String(http.StatusNotFound, "flag not found")
		return
	}
	if err != nil {
		c.String(http.StatusInternalServerError, err.Error())
		return
	}
	c.JSON(http.StatusOK, flag)
}

func (h *FlagsController) Create(c *gin.Context) {
	var req model.Flag
	if err := c.ShouldBindJSON(&req); err != nil {
		c.String(http.StatusBadRequest, "invalid json body")
		return
	}
	flag, err := h.svc.Create(c.Request.Context(), req)
	if err != nil {
		if errors.Is(err, service.ErrInvalidName) {
			c.String(http.StatusBadRequest, err.Error())
			return
		}
		c.String(http.StatusInternalServerError, err.Error())
		return
	}
	c.JSON(http.StatusCreated, flag)
}

func (h *FlagsController) Toggle(c *gin.Context) {
	flag, err := h.svc.Toggle(c.Request.Context(), c.Param("name"))
	if errors.Is(err, service.ErrFlagNotFound) {
		c.String(http.StatusNotFound, "flag not found")
		return
	}
	if err != nil {
		c.String(http.StatusInternalServerError, err.Error())
		return
	}
	c.JSON(http.StatusOK, flag)
}

func (h *FlagsController) Delete(c *gin.Context) {
	if err := h.svc.Delete(c.Request.Context(), c.Param("name")); err != nil {
		if errors.Is(err, service.ErrFlagNotFound) {
			c.String(http.StatusNotFound, "flag not found")
			return
		}
		c.String(http.StatusInternalServerError, err.Error())
		return
	}
	c.Status(http.StatusNoContent)
}
