package repository

import (
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/lib/pq"
	_ "github.com/lib/pq"
	"github.com/obadoraibu/go-auth/internal/config"
	"github.com/obadoraibu/go-auth/internal/domain"
)

type PostgresRepository struct {
	db     *sql.DB
	config *config.PostgresRepositoryConfig
}

func NewPostgresRepository(config *config.PostgresRepositoryConfig) (*PostgresRepository, error) {
	db, err := sql.Open("postgres", fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable",
		config.User, config.Password, config.Host, config.Port, config.Name))
	if err != nil {
		return nil, err
	}
	if err := db.Ping(); err != nil {
		return nil, err
	}

	return &PostgresRepository{
		db:     db,
		config: config,
	}, nil
}

func (r *PostgresRepository) Close() error {
	err := r.db.Close()
	if err != nil {
		return err
	}
	return nil
}

func (r *Repository) FindUserByEmail(email string) (*domain.User, error) {
	u := &domain.User{Email: email}
	if err := r.Postgres.db.QueryRow("SELECT * FROM \"users\"  WHERE email = $1",
		email).Scan(&u.Id, &u.Username, &u.Email, &u.PasswordHash, &u.Role, &u.Status, &u.InviteToken, &u.InviteTokenExp); err != nil {
		if err == sql.ErrNoRows {
			return nil, domain.ErrUserNotFound
		}
		return nil, err
	}
	return u, nil
}

func (r *Repository) GetUserForUpdate(userID int64) (*domain.User, error) {
	query := `
		SELECT id, username, email, role, status, password_hash, invite_token, invite_token_expires_at
		FROM users 
		WHERE id = $1
	`

	row := r.Postgres.db.QueryRow(query, userID)

	u := &domain.User{}
	err := row.Scan(
		&u.Id,
		&u.Username,
		&u.Email,
		&u.Role,
		&u.Status,
		&u.PasswordHash,
		&u.InviteToken,
		&u.InviteTokenExp,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, domain.ErrUserNotFound
		}
		return nil, err
	}

	return u, nil
}

func (r *Repository) GetUserByID(userID int64) (*domain.User, error) {
	query := `
		SELECT id, username, email, role, status
		FROM users 
		WHERE id = $1
	`

	row := r.Postgres.db.QueryRow(query, userID)

	u := &domain.User{}
	err := row.Scan(
		&u.Id,
		&u.Username,
		&u.Email,
		&u.Role,
		&u.Status,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, domain.ErrUserNotFound
		}
		return nil, err
	}

	return u, nil
}

func (r *Repository) CreateUserInvite(u *domain.User) (*domain.User, error) {
	query := `
		INSERT INTO "users" 
			(username, email, role, status, invite_token, invite_token_expires_at)
		VALUES 
			($1, $2, $3, $4, $5, $6)
		RETURNING id
	`

	err := r.Postgres.db.QueryRow(query,
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

func (r *Repository) CompleteInvite(code string, passwordHash string) error {
	tx, err := r.Postgres.db.Begin()
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
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return domain.ErrWrongInviteCode
		}
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

	rows, err := r.Postgres.db.Query(query, role, status, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []*domain.User
	for rows.Next() {
		var u domain.User
		if err := rows.Scan(&u.Id, &u.Username, &u.Email, &u.Role, &u.Status); err != nil {
			return nil, err
		}
		users = append(users, &u)
	}

	return users, nil
}

func (r *Repository) UpdateUserByID(user *domain.User) error {
	query := `
		UPDATE users
		SET 
			username = $1,
			role = $2,
			status = $3,
			password_hash = $4,
			invite_token = $5,
			invite_token_expires_at = $6
		WHERE id = $7
	`

	_, err := r.Postgres.db.Exec(query,
		user.Username,
		user.Role,
		user.Status,
		user.PasswordHash,
		user.InviteToken,
		user.InviteTokenExp,
		user.Id,
	)

	return err
}

func (r *Repository) SaveResetToken(userID int64, resetToken string, expiresAt time.Time) error {
	query := `
		UPDATE users
		SET invite_token = $1,
		    invite_token_expires_at = $2
		WHERE id = $3
	`
	_, err := r.Postgres.db.Exec(query, resetToken, expiresAt, userID)
	return err
}

func (r *Repository) FindUserByResetToken(token string) (*domain.User, error) {
	query := `
		SELECT id, username, email, password_hash, role, status, invite_token, invite_token_expires_at
		FROM users
		WHERE invite_token = $1
	`

	user := &domain.User{}
	err := r.Postgres.db.QueryRow(query, token).Scan(
		&user.Id,
		&user.Username,
		&user.Email,
		&user.PasswordHash,
		&user.Role,
		&user.Status,
		&user.InviteToken,
		&user.InviteTokenExp,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, domain.ErrResetTokenNotFound
		}
		return nil, err
	}

	return user, nil
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
) (int64, error) {
	query := `
		INSERT INTO mri_scans (
			user_id, object_name, original_filename, content_type, size,
			patient_name, patient_gender, patient_age, scan_date
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
		RETURNING id
	`

	var scanID int64
	err := r.Postgres.db.QueryRow(query,
		userID,
		objectName,
		originalFilename,
		contentType,
		size,
		patientName,
		patientGender,
		patientAge,
		scanDate,
	).Scan(&scanID)

	if err != nil {
		return 0, fmt.Errorf("failed to insert scan metadata: %w", err)
	}

	return scanID, nil
}

func (r *Repository) GetScansByFilters(userID int64, filter *domain.ScanFilter) ([]*domain.MRIScan, error) {
	query := `
		SELECT id, user_id, patient_name, patient_gender, patient_age, scan_date, created_at, status
		FROM mri_scans
		WHERE user_id = $1
	`

	args := []interface{}{userID}
	argIdx := 2

	if filter.ScanID != nil {
		query += fmt.Sprintf(" AND id = $%d", argIdx)
		args = append(args, *filter.ScanID)
		argIdx++
	}
	if filter.UploadedFrom != nil {
		query += fmt.Sprintf(" AND created_at >= $%d", argIdx)
		args = append(args, *filter.UploadedFrom)
		argIdx++
	}
	if filter.UploadedTo != nil {
		query += fmt.Sprintf(" AND created_at <= $%d", argIdx)
		args = append(args, *filter.UploadedTo)
		argIdx++
	}

	query += " ORDER BY created_at DESC"

	rows, err := r.Postgres.db.Query(query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var scans []*domain.MRIScan
	for rows.Next() {
		var s domain.MRIScan
		err := rows.Scan(
			&s.ID,
			&s.UserID,
			&s.PatientName,
			&s.PatientGender,
			&s.PatientAge,
			&s.ScanDate,
			&s.CreatedAt,
			&s.Status,
		)
		if err != nil {
			return nil, err
		}
		scans = append(scans, &s)
	}

	return scans, nil
}

func (r *Repository) GetScanDetail(userID, scanID int64) (*domain.MRIScanDetail, error) {
	query := `
		SELECT s.id, s.user_id, s.patient_name, s.patient_gender, s.patient_age,
		       s.scan_date, s.object_name, s.original_filename, s.content_type, s.size,
		       s.created_at, s.status,
		       t.diagnosis, t.confidence, t.gradcam_url, t.completed_at
		FROM mri_scans s
		LEFT JOIN mri_analysis_results t ON s.id = t.scan_id
		WHERE s.user_id = $1 AND s.id = $2
	`

	scan := &domain.MRIScanDetail{}
	err := r.Postgres.db.QueryRow(query, userID, scanID).Scan(
		&scan.ID, &scan.UserID, &scan.PatientName, &scan.PatientGender,
		&scan.PatientAge, &scan.ScanDate, &scan.ObjectName, &scan.OriginalName,
		&scan.ContentType, &scan.Size, &scan.CreatedAt, &scan.Status,
		&scan.Diagnosis, &scan.Confidence, &scan.GradCAMURL, &scan.CompletedAt,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("scan not found")
		}
		return nil, err
	}

	return scan, nil
}

func (r *Repository) UpdateUserPassword(userID int64, hash string) error {
	query := `
		UPDATE users
		SET password_hash = $1,
		    invite_token = NULL,
		    invite_token_expires_at = NULL
		WHERE id = $2
	`

	_, err := r.Postgres.db.Exec(query, hash, userID)
	return err
}
