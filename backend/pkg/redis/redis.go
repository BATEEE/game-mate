package redis

import (
	"context"
	"fmt"
	"time"

	"github.com/BATEEE/game-mate/backend/pkg/config"
	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"
)

func NewRedisClient(ctx context.Context, cfg *config.RedisConfig, log *zap.Logger) (*redis.Client, error) {
	client := redis.NewClient(&redis.Options{
		Addr:     cfg.Addr(),
		Password: cfg.Password,
		DB:       cfg.DB,
	})

	pingCtx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()
	if err := client.Ping(pingCtx).Err(); err != nil {
		client.Close()
		return nil, fmt.Errorf("không thể kết nối tới redis: %w", err)
	}
	
	log.Info("Kết nối Redis thành công",
		zap.String("addr", cfg.Addr()),
		zap.Int("db", cfg.DB),
	)
	return client, nil
}
