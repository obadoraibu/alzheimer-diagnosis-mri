package repository

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/lib/pq"
	_ "github.com/lib/pq"
	"github.com/obadoraibu/go-auth/internal/config"
	"github.com/obadoraibu/go-auth/internal/domain"
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

func (r *Repository) CreateUserInvite(u *domain.User) (*domain.User, error) {
	query := `
		INSERT INTO "users" 
			(username, email, role, status, invite_token, invite_token_expires_at)
		VALUES 
			($1, $2, $3, $4, $5, $6)
		RETURNING id
	`

	err := r.Users.db.QueryRow(query,
		u.Username,
		u.Email,
		u.Role,
		u.Status,
		u.InviteToken.String,
		u.InviteTokenExp.Time,
	).Scan(&u.Id)

	if err != nil {
		if pqErr, ok := err.(*pq.Error); ok && pqErr.Code == "23505" {
			return nil, domain.ErrUserAlreadyExists
		}
		return nil, err
	}

	return u, nil
}

func (r *Repository) FindUserByEmail(email string) (*domain.User, error) {
	u := &domain.User{Email: email}
	if err := r.Users.db.QueryRow("SELECT * FROM \"users\"  WHERE email = $1",
		email).Scan(&u.Id, &u.Username, &u.Email, &u.PasswordHash, &u.Role, &u.Status, &u.InviteToken, &u.InviteTokenExp); err != nil {
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

func (r *Repository) CompleteInvite(code string, passwordHash string) error {
	tx, err := r.Users.db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	var userID int
	var currentStatus string
	var expiresAt time.Time

	query := `
		SELECT id, status, invite_token_expires_at 
		FROM users 
		WHERE invite_token = $1
	`
	err = tx.QueryRow(query, code).Scan(&userID, &currentStatus, &expiresAt)
	if err == sql.ErrNoRows {
		return domain.ErrWrongInviteCode
	}
	if err != nil {
		return err
	}

	if currentStatus != "invited" {
		return domain.ErrInviteAlreadyUsed
	}

	if time.Now().After(expiresAt) {
		return domain.ErrInviteExpired
	}

	updateQuery := `
		UPDATE users 
		SET password_hash = $1,
		    status = 'active',
		    invite_token = NULL,
		    invite_token_expires_at = NULL
		WHERE id = $2
	`
	_, err = tx.Exec(updateQuery, passwordHash, userID)
	if err != nil {
		return err
	}

	return tx.Commit()
}
