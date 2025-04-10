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

type TokenRedisRepository struct {
	client          *redis.Client
	config          *config.TokenRepositoryConfig
	refreshTokenTTL string
}

func NewTokenRepository(config *config.TokenRepositoryConfig) (*TokenRedisRepository, error) {
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

	return &TokenRedisRepository{
		client: client,
		config: config,
	}, nil
}

func (r *TokenRedisRepository) Close() error {
	err := r.client.Close()
	if err != nil {
		return err
	}
	return nil
}

// new
func (r *Repository) AddToken(fingerprint, refresh, email, role string) error {
	type TokenData struct {
		Email string `json:"email"`
		Role  string `json:"role"`
	}

	ctx := context.Background()
	ttl := time.Hour * 24 * 60

	key := fmt.Sprintf("%s:%s", fingerprint, refresh)

	data := TokenData{
		Email: email,
		Role:  role,
	}

	value, err := json.Marshal(data)
	if err != nil {
		return err
	}

	err = r.Tokens.client.Set(ctx, key, value, ttl).Err()
	if err != nil {
		return err
	}

	return nil
}

func (r *Repository) FindAndDeleteRefreshToken(refresh, fingerprint string) (string, error) {
	key := fmt.Sprintf("%s:%s", fingerprint, refresh)

	exists, err := r.Tokens.client.Exists(context.Background(), key).Result()
	if err != nil {
		return "", err
	}

	if exists == 0 {
		err := errors.New("key does not exist")
		return "", err
	} else {
		value, err := r.Tokens.client.Get(context.Background(), key).Result()
		if err != nil {
			return "", err
		}
		err = r.Tokens.client.Del(context.Background(), key).Err()
		if err != nil {
			panic(err)
		}
		return value, nil
	}
}

func (r *Repository) DeleteToken(u *domain.User) error { return nil }
