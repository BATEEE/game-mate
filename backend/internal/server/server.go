package server

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/BATEEE/game-mate/backend/internal/health"
	"github.com/BATEEE/game-mate/backend/pkg/config"
	"go.uber.org/zap"
)

// Server quản lý vòng đời và Graceful Shutdown của HTTP Server (SRP).
type Server struct {
	httpServer *http.Server
	log        *zap.Logger
}

func New(cfg *config.Config, log *zap.Logger, healthHandler *health.Handler) *Server {
	router := NewRouter(RouterConfig{
		Mode:          cfg.Server.Mode,
		Logger:        log,
		HealthHandler: healthHandler,
	})

	addr := fmt.Sprintf(":%s", cfg.Server.Port)

	httpSrv := &http.Server{
		Addr:         addr,
		Handler:      router,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	return &Server{
		httpServer: httpSrv,
		log:        log,
	}
}

// Run khởi động HTTP server và lắng nghe tín hiệu tắt server an toàn.
func (s *Server) Run() error {
	serverErrors := make(chan error, 1)

	go func() {
		s.log.Info("HTTP Server đang lắng nghe", zap.String("addr", s.httpServer.Addr))
		if err := s.httpServer.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			serverErrors <- err
		}
	}()

	shutdown := make(chan os.Signal, 1)
	signal.Notify(shutdown, os.Interrupt, syscall.SIGTERM, syscall.SIGINT)

	select {
	case err := <-serverErrors:
		return fmt.Errorf("lỗi khởi động server: %w", err)

	case sig := <-shutdown:
		s.log.Info("Nhận tín hiệu dừng, bắt đầu Graceful Shutdown...", zap.String("signal", sig.String()))

		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		if err := s.httpServer.Shutdown(ctx); err != nil {
			s.httpServer.Close()
			return fmt.Errorf("buộc phải dừng server do lỗi shutdown: %w", err)
		}

		s.log.Info("HTTP Server đã dừng an toàn")
	}

	return nil
}
