package database

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"strconv"
	"time"

	"github.com/ericp/chronos-bot-reminder/internal/config"
	"github.com/redis/go-redis/v9"
)

var RedisClient *redis.Client
var redisCtx = context.Background()

// InitializeRedis sets up the Redis connection
func InitializeRedis() error {
	cfg := config.Load()

	db, err := strconv.Atoi(cfg.RedisDB)
	if err != nil {
		return fmt.Errorf("[REDIS] - Invalid Redis DB number: %w", err)
	}

	RedisClient = redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%s", cfg.RedisHost, cfg.RedisPort),
		Password: cfg.RedisPassword,
		DB:       db,
	})

	// Test the connection
	_, err = RedisClient.Ping(redisCtx).Result()
	if err != nil {
		return fmt.Errorf("[REDIS] - Failed to connect to Redis: %w", err)
	}

	log.Println("[REDIS] - ✅ Connection established")
	return nil
}

// GetRedisClient returns the Redis client instance
func GetRedisClient() *redis.Client {
	return RedisClient
}

// Set stores a key-value pair in Redis with optional expiration
func SetCache(key string, value interface{}, expiration time.Duration) error {
	jsonValue, err := json.Marshal(value)
	if err != nil {
		return fmt.Errorf("failed to marshal value: %w", err)
	}

	return RedisClient.Set(redisCtx, key, jsonValue, expiration).Err()
}

// Get retrieves a value from Redis and unmarshals it into the provided interface
func GetCache(key string, dest interface{}) error {
	val, err := RedisClient.Get(redisCtx, key).Result()
	if err != nil {
		if err == redis.Nil {
			return fmt.Errorf("key %s not found", key)
		}
		return fmt.Errorf("failed to get value: %w", err)
	}

	return json.Unmarshal([]byte(val), dest)
}

// Delete removes a key from Redis
func DeleteCache(key string) error {
	return RedisClient.Del(redisCtx, key).Err()
}

// Exists checks if a key exists in Redis
func ExistsCache(key string) (bool, error) {
	result, err := RedisClient.Exists(redisCtx, key).Result()
	return result > 0, err
}

// SetString stores a string value in Redis
func SetStringCache(key, value string, expiration time.Duration) error {
	return RedisClient.Set(redisCtx, key, value, expiration).Err()
}

// GetString retrieves a string value from Redis
func GetStringCache(key string) (string, error) {
	return RedisClient.Get(redisCtx, key).Result()
}

// SetHash stores a hash map in Redis
func SetHashCache(key string, fields map[string]interface{}) error {
	return RedisClient.HMSet(redisCtx, key, fields).Err()
}

// GetHash retrieves a hash map from Redis
func GetHashCache(key string) (map[string]string, error) {
	return RedisClient.HGetAll(redisCtx, key).Result()
}

// SetExpire sets expiration time for an existing key
func SetExpireCache(key string, expiration time.Duration) error {
	return RedisClient.Expire(redisCtx, key, expiration).Err()
}

// CloseRedis closes the Redis connection
func CloseRedis() error {
	return RedisClient.Close()
}