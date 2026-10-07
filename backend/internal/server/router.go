package server

import (
	"github.com/BATEEE/game-mate/backend/internal/health"
	"github.com/BATEEE/game-mate/backend/internal/middleware"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// RouterConfig chứa các dependency cần thiết để cấu hình router.
type RouterConfig struct {
	Mode          string
	Logger        *zap.Logger
	HealthHandler *health.Handler
}

// NewRouter cấu hình và trả về Gin Engine (SRP).
func NewRouter(cfg RouterConfig) *gin.Engine {
	if cfg.Mode == "release" {
		gin.SetMode(gin.ReleaseMode)
	} else {
		gin.SetMode(gin.DebugMode)
	}

	r := gin.New()

	// Middlewares
	r.Use(gin.Recovery())
	r.Use(middleware.ZapLogger(cfg.Logger))
	r.Use(middleware.CORS())

	// Root API Group
	api := r.Group("/api")
	{
		cfg.HealthHandler.RegisterRoutes(api)
	}

	return r
}
