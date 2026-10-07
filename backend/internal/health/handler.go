package health

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// Handler chịu trách nhiệm duy nhất cho tầng HTTP Transport của Health (SRP).
type Handler struct {
	svc Service
}

func NewHandler(svc Service) *Handler {
	return &Handler{
		svc: svc,
	}
}

func (h *Handler) RegisterRoutes(rg *gin.RouterGroup) {
	rg.GET("/health", h.Check)
}

func (h *Handler) Check(c *gin.Context) {
	result := h.svc.Check(c.Request.Context())

	statusCode := http.StatusOK
	if result.Status == StatusDown {
		statusCode = http.StatusServiceUnavailable
	}

	c.JSON(statusCode, result)
}
