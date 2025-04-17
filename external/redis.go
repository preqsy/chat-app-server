package external

import (
	"context"

	"github.com/redis/go-redis/v9"
	"github.com/sirupsen/logrus"
)

type RedisService struct {
	rdb    *redis.Client
	ctx    context.Context
	logger *logrus.Logger
}

func InitRedis(ctx context.Context, logger *logrus.Logger, redisUrl string) (*RedisService, error) {
	logger.Info("Connecting to redis.....")

	opt, err := redis.ParseURL(redisUrl)
	if err != nil {
		logger.Errorf("Redis ParseURL failed: %v", err)
		return nil, err
	}

	rdb := redis.NewClient(opt)

	if err := rdb.Ping(ctx).Err(); err != nil {
		logger.Errorf("Redis connection failed: %v", err)
		return nil, err
	}

	logger.Info("Connection with Redis established")
	return &RedisService{rdb: rdb, ctx: ctx, logger: logger}, nil
}

func (r *RedisService) PublishMessage(channel, message string) error {
	return r.rdb.Publish(r.ctx, channel, message).Err()
}

func (r *RedisService) SubscribeToChannel(channel string) *redis.PubSub {
	return r.rdb.Subscribe(r.ctx, channel)
}
