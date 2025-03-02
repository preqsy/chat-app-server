package external

import (
	"context"

	"github.com/redis/go-redis/v9"
)

type RedisService struct {
	rdb *redis.Client
	ctx context.Context
}

func InitRedis(ctx context.Context) (*RedisService, error) {
	rdb := redis.NewClient(&redis.Options{Addr: "localhost:6379"})

	if err := rdb.Ping(ctx).Err(); err != nil {
		return nil, err
	}
	return &RedisService{rdb: rdb, ctx: ctx}, nil
}

func (r *RedisService) PublishMessage(channel, message string) error {
	return r.rdb.Publish(r.ctx, channel, message).Err()
}

func (r *RedisService) SubscribeToChannel(channel string) *redis.PubSub {
	return r.rdb.Subscribe(r.ctx, channel)
}
