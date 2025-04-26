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

// new
func (r *Repository) AddToken(fingerprint, refresh string, user_id int64, role string) error {
	type TokenData struct {
		Id   string `json:"user_id"`
		Role string `json:"role"`
	}

	ctx := context.Background()
	ttl := time.Hour * 24 * 60

	key := fmt.Sprintf("%s:%s", fingerprint, refresh)

	data := TokenData{
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
	key := fmt.Sprintf("%s:%s", fingerprint, refresh)

	exists, err := r.Redis.client.Exists(context.Background(), key).Result()
	if err != nil {
		return "", err
	}

	if exists == 0 {
		err := errors.New("key does not exist")
		return "", err
	} else {
		value, err := r.Redis.client.Get(context.Background(), key).Result()
		if err != nil {
			return "", err
		}
		err = r.Redis.client.Del(context.Background(), key).Err()
		if err != nil {
			panic(err)
		}
		return value, nil
	}
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
