package redis

import (
	"context"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
	"github.com/ydcloud-dy/leaf-api/config"
)

var (
	Client *redis.Client
	ctx    = context.Background()
)

// InitRedis 初始化 Redis 客户端
func InitRedis() error {
	cfg := config.AppConfig.Redis

	Client = redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%d", cfg.Host, cfg.Port),
		Password: cfg.Password,
		DB:       cfg.DB,
		PoolSize: cfg.PoolSize,
	})

	// 测试连接
	if err := Client.Ping(ctx).Err(); err != nil {
		return fmt.Errorf("failed to connect to redis: %w", err)
	}

	return nil
}

// Close 关闭 Redis 连接
func Close() error {
	if Client != nil {
		return Client.Close()
	}
	return nil
}

// GetClient 获取 Redis 客户端
func GetClient() *redis.Client {
	return Client
}

// GetContext 获取默认 context
func GetContext() context.Context {
	return ctx
}

// SetWithExpire 设置带过期时间的 key
func SetWithExpire(key string, value interface{}, expiration time.Duration) error {
	return Client.Set(ctx, key, value, expiration).Err()
}

// Get 获取 key 的值
func Get(key string) (string, error) {
	return Client.Get(ctx, key).Result()
}

// Del 删除 key
func Del(key string) error {
	return Client.Del(ctx, key).Err()
}

// Exists 检查 key 是否存在
func Exists(key string) (bool, error) {
	result, err := Client.Exists(ctx, key).Result()
	return result > 0, err
}

// Keys 根据模式获取所有匹配的 key
func Keys(pattern string) ([]string, error) {
	return Client.Keys(ctx, pattern).Result()
}

// Expire 设置 key 的过期时间
func Expire(key string, expiration time.Duration) error {
	return Client.Expire(ctx, key, expiration).Err()
}

// Incr 递增计数器
func Incr(key string) (int64, error) {
	return Client.Incr(ctx, key).Result()
}

// Decr 递减计数器
func Decr(key string) (int64, error) {
	return Client.Decr(ctx, key).Result()
}

// GetInt 获取整数值
func GetInt(key string) (int64, error) {
	val, err := Client.Get(ctx, key).Int64()
	if err == redis.Nil {
		return 0, nil
	}
	return val, err
}

// SetInt 设置整数值
func SetInt(key string, value int64, expiration time.Duration) error {
	return Client.Set(ctx, key, value, expiration).Err()
}
