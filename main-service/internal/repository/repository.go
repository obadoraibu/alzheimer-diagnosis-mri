package repository

import (
	"github.com/obadoraibu/go-auth/internal/config"
)

type Repository struct {
	Users  *UserPostgresRepository
	Tokens *TokenRedisRepository
	MinIO  *MinIOClient
	config *config.DatabaseConfig
}

func (r *Repository) Close() error {
	err := r.Users.Close()
	if err != nil {
		return err
	}

	err = r.Tokens.Close()
	if err != nil {
		return err
	}
	return nil
}

func NewRepository(config *config.DatabaseConfig) (*Repository, error) {
	var repo Repository

	repo.config = config

	userRepository, err := NewUserRepository(config.UserRepositoryConfig)
	if err != nil {
		return nil, err
	}
	repo.Users = userRepository

	tokenRepository, err := NewTokenRepository(config.TokenRepositoryConfig)
	if err != nil {
		return nil, err
	}
	repo.Tokens = tokenRepository

	minioClient, err := NewMinIOClient(config.MinIORepositoryConfig)
	if err != nil {
		return nil, err
	}
	repo.MinIO = minioClient

	return &repo, nil
}
