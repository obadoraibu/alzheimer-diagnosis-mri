package repository

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/obadoraibu/go-auth/internal/config"
	"github.com/obadoraibu/go-auth/internal/domain"
	"github.com/redis/go-redis/v9"
	"github.com/sirupsen/logrus"
)

type RedisRepository struct {
	client          *redis.Client
	config          *config.RedisRepositoryConfig
	refreshTokenTTL string
}

func NewRedisRepository(config *config.RedisRepositoryConfig) (*RedisRepository, error) {
	redisHost := config.Host
	redisPort := config.Port
	redisPassword := config.Password

	client := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%s", redisHost, redisPort),
		Password: redisPassword,
		DB:       0,
	})

	err := client.Ping(context.Background()).Err()
	if err != nil {
		return nil, err
	}

	logrus.Print("connected to redis, %s", err)

	return &RedisRepository{
		client: client,
		config: config,
	}, nil
}

func (r *RedisRepository) Close() error {
	err := r.client.Close()
	if err != nil {
		return err
	}
	return nil
}

func (r *Repository) AddToken(fingerprint, refresh string, user_id int64, role string) error {
	ctx := context.Background()
	ttl := time.Hour * 24 * 60

	key := fmt.Sprintf("refresh:%s:%s", fingerprint, refresh)

	data := domain.TokenData{
		Id:   fmt.Sprintf("%d", user_id),
		Role: role,
	}

	value, err := json.Marshal(data)
	if err != nil {
		return err
	}

	err = r.Redis.client.Set(ctx, key, value, ttl).Err()
	if err != nil {
		return err
	}

	return nil
}

func (r *Repository) FindAndDeleteRefreshToken(refresh, fingerprint string) (string, error) {
	ctx := context.Background()
	key := fmt.Sprintf("refresh:%s:%s", fingerprint, refresh)

	value, err := r.Redis.client.Get(ctx, key).Result()
	if err != nil {
		if errors.Is(err, redis.Nil) {
			return "", domain.ErrTokenNotFound
		}
		return "", err
	}

	if err := r.Redis.client.Del(ctx, key).Err(); err != nil {
		return "", domain.ErrInternal(err)
	}

	return value, nil
}

func (r *Repository) DeleteToken(u *domain.User) error { return nil }

func (r *Repository) EnqueueScanTask(scanID int64, objectName string) error {
	ctx := context.Background()

	task := struct {
		ScanID     int64  `json:"scan_id"`
		ObjectName string `json:"object_name"`
		CreatedAt  int64  `json:"created_at"`
	}{
		ScanID:     scanID,
		ObjectName: objectName,
		CreatedAt:  time.Now().Unix(),
	}

	data, err := json.Marshal(task)
	if err != nil {
		return fmt.Errorf("failed to marshal task: %w", err)
	}

	queueName := "mri_tasks"
	err = r.Redis.client.RPush(ctx, queueName, data).Err()
	if err != nil {
		return fmt.Errorf("failed to enqueue task: %w", err)
	}

	logrus.Infof("Task enqueued for user %d: %s", scanID, objectName)
	return nil
}
