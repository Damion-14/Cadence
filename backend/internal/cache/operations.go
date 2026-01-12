package cache

import (
	"context"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
)

const (
	KeyActiveWorkout   = "active_workout:user:%d"
	KeyUserPRs         = "prs:user:%d"
	KeyWeeklySummary   = "weekly:user:%d:week:%s"
	KeyExerciseProgress = "progress:user:%d:exercise:%s"
)

const (
	TTLActiveWorkout   = 24 * time.Hour
	TTLUserPRs         = 1 * time.Hour
	TTLWeeklySummary   = 7 * 24 * time.Hour
	TTLExerciseProgress = 1 * time.Hour
)

type Cache struct {
	client *redis.Client
}

func NewCache(client *redis.Client) *Cache {
	return &Cache{client: client}
}

func (c *Cache) Get(ctx context.Context, key string) (string, error) {
	return c.client.Get(ctx, key).Result()
}

func (c *Cache) Set(ctx context.Context, key string, value interface{}, ttl time.Duration) error {
	return c.client.Set(ctx, key, value, ttl).Err()
}

func (c *Cache) Delete(ctx context.Context, keys ...string) error {
	if len(keys) == 0 {
		return nil
	}
	return c.client.Del(ctx, keys...).Err()
}

func (c *Cache) DeletePattern(ctx context.Context, pattern string) error {
	iter := c.client.Scan(ctx, 0, pattern, 0).Iterator()
	var keys []string

	for iter.Next(ctx) {
		keys = append(keys, iter.Val())
	}

	if err := iter.Err(); err != nil {
		return err
	}

	if len(keys) > 0 {
		return c.client.Del(ctx, keys...).Err()
	}

	return nil
}

func (c *Cache) Exists(ctx context.Context, key string) (bool, error) {
	count, err := c.client.Exists(ctx, key).Result()
	if err != nil {
		return false, err
	}
	return count > 0, nil
}

func GetActiveWorkoutKey(userID int) string {
	return fmt.Sprintf(KeyActiveWorkout, userID)
}

func GetUserPRsKey(userID int) string {
	return fmt.Sprintf(KeyUserPRs, userID)
}

func GetWeeklySummaryKey(userID int, week string) string {
	return fmt.Sprintf(KeyWeeklySummary, userID, week)
}

func GetExerciseProgressKey(userID int, exerciseName string) string {
	return fmt.Sprintf(KeyExerciseProgress, userID, exerciseName)
}
