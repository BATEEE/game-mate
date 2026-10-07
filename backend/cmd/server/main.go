package main

import (
	"context"
	"fmt"
	"time"

	"github.com/BATEEE/game-mate/backend/internal/health"
	"github.com/BATEEE/game-mate/backend/internal/server"
	"github.com/BATEEE/game-mate/backend/pkg/config"
	"github.com/BATEEE/game-mate/backend/pkg/database"
	"github.com/BATEEE/game-mate/backend/pkg/logger"
	"github.com/BATEEE/game-mate/backend/pkg/redis"
	"go.uber.org/zap"
)

func main() {
	// 1. Đọc cấu hình
	cfg, err := config.Load()
	if err != nil {
		panic(fmt.Sprintf("Failed to load config: %v", err))
	}

	// 2. Khởi tạo logger
	log, err := logger.New(cfg.Server.Mode)
	if err != nil {
		panic(fmt.Sprintf("Failed to init logger: %v", err))
	}
	defer log.Sync()

	log.Info("GameMate Backend đang khởi động...",
		zap.String("port", cfg.Server.Port),
		zap.String("mode", cfg.Server.Mode),
	)

	// Context dùng cho khởi tạo hạ tầng ban đầu
	initCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// 3. Kết nối PostgreSQL
	dbPool, err := database.NewPostgresPool(initCtx, &cfg.Database, log)
	if err != nil {
		log.Fatal("Không thể kết nối PostgreSQL", zap.Error(err))
	}
	defer dbPool.Close()

	// 4. Kết nối Redis
	redisClient, err := redis.NewRedisClient(initCtx, &cfg.Redis, log)
	if err != nil {
		log.Fatal("Không thể kết nối Redis", zap.Error(err))
	}
	defer redisClient.Close()

	// 5. Khởi tạo Health Check (SOLID: OCP + DIP)
	checkers := []health.Checker{
		health.NewPostgresChecker(dbPool),
		health.NewRedisChecker(redisClient),
	}
	healthService := health.NewService(checkers)
	healthHandler := health.NewHandler(healthService)

	// 6. Khởi tạo và khởi chạy HTTP Server (Graceful Shutdown)
	srv := server.New(cfg, log, healthHandler)
	if err := srv.Run(); err != nil {
		log.Fatal("Server dừng bất thường", zap.Error(err))
	}
}