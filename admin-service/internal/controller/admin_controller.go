package controller

import (
	"io"
	"net/http"

	"github.com/gin-gonic/gin"

	"featureflags/admin-service/internal/service"
)

// AdminController exposes dashboard API routes.
type AdminController struct {
	svc *service.AdminService
}

func NewAdminController(svc *service.AdminService) *AdminController {
	return &AdminController{svc: svc}
}

func (h *AdminController) RegisterRoutes(r *gin.Engine) {
	r.GET("/health", h.Health)

	api := r.Group("/api")
	api.GET("/flags", h.List)
	api.POST("/flags", h.Create)
	api.PATCH("/flags/:name/toggle", h.Toggle)
	api.DELETE("/flags/:name", h.Delete)
}

func (h *AdminController) Health(c *gin.Context) {
	c.String(http.StatusOK, "ok")
}

func (h *AdminController) List(c *gin.Context) {
	status, body, err := h.svc.List()
	if err != nil {
		c.String(http.StatusBadGateway, err.Error())
		return
	}
	c.Data(status, "application/json", body)
}

func (h *AdminController) Create(c *gin.Context) {
	body, err := io.ReadAll(c.Request.Body)
	if err != nil {
		c.String(http.StatusBadRequest, "failed reading body")
		return
	}
	status, resp, err := h.svc.CreateFlag(body)
	if err != nil {
		c.String(http.StatusBadGateway, err.Error())
		return
	}
	c.Data(status, "application/json", resp)
}

func (h *AdminController) Toggle(c *gin.Context) {
	status, resp, err := h.svc.ToggleFlag(c.Param("name"))
	if err != nil {
		c.String(http.StatusBadGateway, err.Error())
		return
	}
	c.Data(status, "application/json", resp)
}

func (h *AdminController) Delete(c *gin.Context) {
	status, resp, err := h.svc.DeleteFlag(c.Param("name"))
	if err != nil {
		c.String(http.StatusBadGateway, err.Error())
		return
	}
	if len(resp) > 0 {
		c.Data(status, "application/json", resp)
	} else {
		c.Status(status)
	}
}
