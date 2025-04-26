package repository

import (
	"github.com/obadoraibu/go-auth/internal/config"
)

type Repository struct {
	Postgres *PostgresRepository
	Redis    *RedisRepository
	MinIO    *MinIOClient
	config   *config.DatabaseConfig
}

func (r *Repository) Close() error {
	err := r.Postgres.Close()
	if err != nil {
		return err
	}

	err = r.Redis.Close()
	if err != nil {
		return err
	}
	return nil
}

func NewRepository(config *config.DatabaseConfig) (*Repository, error) {
	var repo Repository

	repo.config = config

	userRepository, err := NewPostgresRepository(config.PostgresRepositoryConfig)
	if err != nil {
		return nil, err
	}
	repo.Postgres = userRepository

	tokenRepository, err := NewRedisRepository(config.RedisRepositoryConfig)
	if err != nil {
		return nil, err
	}
	repo.Redis = tokenRepository

	minioClient, err := NewMinIOClient(config.MinIORepositoryConfig)
	if err != nil {
		return nil, err
	}
	repo.MinIO = minioClient

	return &repo, nil
}
