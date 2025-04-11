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

func (r *Repository) GetUsersFiltered(role, status string, limit, offset int) ([]*domain.User, error) {
	query := `
		SELECT id, username, email, role, status 
		FROM users
		WHERE ($1 = '' OR role = $1)
		  AND ($2 = '' OR status = $2)
		ORDER BY id
		LIMIT $3 OFFSET $4;
	`

	rows, err := r.Users.db.Query(query, role, status, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []*domain.User
	for rows.Next() {
		var u domain.User
		err := rows.Scan(&u.Id, &u.Username, &u.Email, &u.Role, &u.Status)
		if err != nil {
			return nil, err
		}
		users = append(users, &u)
	}

	return users, nil
}

func (r *Repository) GetUserByID(userID int64) (*domain.User, error) {
	query := `
		SELECT id, username, email, role, status 
		FROM users 
		WHERE id = $1
	`

	row := r.Users.db.QueryRow(query, userID)

	var user domain.User
	err := row.Scan(&user.Id, &user.Username, &user.Email, &user.Role, &user.Status)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("user not found")
		}
		return nil, err
	}

	return &user, nil
}

func (r *Repository) UpdateUserByID(user *domain.User) error {
	query := `
		UPDATE users
		SET username = $1,
		    role = $2,
		    status = $3
		WHERE id = $4
	`

	_, err := r.Users.db.Exec(query,
		user.Username,
		user.Role,
		user.Status,
		user.Id,
	)

	return err
}

func (r *Repository) SaveScanMetadata(
	userID int64,
	objectName string,
	originalFilename string,
	contentType string,
	size int64,
	patientName string,
	patientGender string,
	patientAge int,
	scanDate time.Time,
) error {
	query := `
		INSERT INTO mri_scans (
			user_id, object_name, original_filename, content_type, size,
			patient_name, patient_gender, patient_age, scan_date
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
	`

	_, err := r.Users.db.Exec(query,
		userID,
		objectName,
		originalFilename,
		contentType,
		size,
		patientName,
		patientGender,
		patientAge,
		scanDate,
	)
	if err != nil {
		return fmt.Errorf("failed to insert scan metadata: %w", err)
	}

	return nil
}
