package repository

import (
	"database/sql"
	"fmt"
	"github.com/lib/pq"
	_ "github.com/lib/pq"
	"github.com/obadoraibu/go-auth/internal/config"
	"github.com/obadoraibu/go-auth/internal/domain"
	"time"
)

type UserPostgresRepository struct {
	db     *sql.DB
	config *config.UserRepositoryConfig
}

func NewUserRepository(config *config.UserRepositoryConfig) (*UserPostgresRepository, error) {
	db, err := sql.Open("postgres", fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable",
		config.User, config.Password, config.Host, config.Port, config.Name))
	if err != nil {
		return nil, err
	}
	if err := db.Ping(); err != nil {
		return nil, err
	}

	return &UserPostgresRepository{
		db:     db,
		config: config,
	}, nil
}

func (r *UserPostgresRepository) Close() error {
	err := r.db.Close()
	if err != nil {
		return err
	}
	return nil
}

func (r *Repository) CreateUserAndEmailConfirmation(u *domain.User, confirmationCode string, expiresAt time.Time) (*domain.User, error) {
	tx, err := r.Users.db.Begin()
	if err != nil {
		return nil, err
	}

	defer tx.Rollback()

	if err := r.Users.db.QueryRow("INSERT INTO \"users\" (username, email, password_hash, is_confirmed) VALUES ($1, $2, $3, $4) RETURNING id",
		u.Name, u.Email, u.PasswordHash, u.IsConfirmed).Scan(&u.Id); err != nil {
		pqErr, ok := err.(*pq.Error)
		if ok {
			if pqErr.Code == "23505" {
				return nil, domain.ErrUserAlreadyExists
			}
		}
		return nil, err
	}

	if _, err := r.Users.db.Exec("INSERT INTO email_confirmations (user_id, code, expires_at) VALUES ($1, $2, $3)",
		u.Id, confirmationCode, expiresAt); err != nil {
		return nil, err
	}

	if err := tx.Commit(); err != nil {
		return nil, err
	}

	return u, nil
}

func (r *Repository) FindUserByEmail(email string) (*domain.User, error) {
	u := &domain.User{Email: email}
	if err := r.Users.db.QueryRow("SELECT * FROM \"users\"  WHERE email = $1",
		email).Scan(&u.Id, &u.Name, &u.Email, &u.PasswordHash, &u.IsConfirmed); err != nil {
		if err == sql.ErrNoRows {
			return nil, domain.ErrWrongEmailOrPassword
		}
		return nil, err
	}
	return u, nil
}

func (r *Repository) ConfirmEmail(code string) error {
	tx, err := r.Users.db.Begin()
	if err != nil {
		return err
	}

	defer tx.Rollback()

	var userID int
	err = tx.QueryRow("SELECT user_id FROM email_confirmations WHERE code = $1 AND expires_at > NOW();", code).Scan(&userID)

	if err == sql.ErrNoRows {
		rollbackErr := tx.Rollback()
		if rollbackErr != nil {
			return domain.ErrWrongEmailConfirmationCode
		}
	} else if err != nil {
		return err
	}

	_, err = tx.Exec("DELETE FROM email_confirmations WHERE code = $1", code)
	if err != nil {
		return err
	}

	_, err = tx.Exec("UPDATE \"users\" SET is_confirmed = $1 WHERE id = $2", true, userID)
	if err != nil {
		return err
	}

	err = tx.Commit()
	if err != nil {
		return err
	}

	return nil
}
