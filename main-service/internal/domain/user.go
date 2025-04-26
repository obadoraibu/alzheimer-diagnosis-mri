package domain

import (
	"database/sql"
	"mime/multipart"
	"time"
)

type User struct {
	Id             int64
	Username       string `json:"username" binding:"required"`
	Email          string `json:"email" binding:"required"`
	PasswordHash   string
	Role           string `json:"role"`
	Status         string `json:"status"`
	InviteToken    sql.NullString
	InviteTokenExp sql.NullTime
}

type CreateUserInvite struct {
	Username string `json:"username" binding:"required"`
	Email    string `json:"email" binding:"required"`
	Role     string `json:"role" binding:"required,oneof=admin doctor"`
}

type UserSignUpInput struct {
	Name     string `json:"name" binding:"required"`
	Email    string `json:"email" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type UserResponse struct {
	ID       int64  `json:"id"`
	Username string `json:"username"`
	Email    string `json:"email"`
	Role     string `json:"role"`
	Status   string `json:"status"`
}

type UpdateUserInput struct {
	Username string `json:"username,omitempty"`
	Role     string `json:"role,omitempty"`
	Status   string `json:"status,omitempty"`
}

// [x]
type UserSignInInput struct {
	Email       string `json:"email" binding:"required"`
	Password    string `json:"password" binding:"required"`
	Fingerprint string `json:"fingerprint" binding:"required"`
}

// [x]
type UserSignInResponse struct {
	AccessToken  string `json:"access" binding:"required"`
	RefreshToken string `json:"refresh"binding:"required"`
}

// [x]
type CreateAnalysisInput struct {
	UserID        int64
	File          multipart.File
	Header        *multipart.FileHeader
	PatientName   string
	PatientGender string
	PatientAge    int
	ScanDate      time.Time
}

// [x]
type ScanFilter struct {
	ScanID       *int64
	UploadedFrom *time.Time
	UploadedTo   *time.Time
}

type MRIScan struct {
	ID            int64
	UserID        int64
	PatientName   string
	PatientGender string
	PatientAge    int
	ScanDate      time.Time
	CreatedAt     time.Time
	Status        string
}

type MRIScanDetail struct {
	ID            int64     `json:"id"`
	UserID        int64     `json:"user_id"`
	PatientName   string    `json:"patient_name"`
	PatientGender string    `json:"patient_gender"`
	PatientAge    int       `json:"patient_age"`
	ScanDate      time.Time `json:"scan_date"`
	ObjectName    string    `json:"object_name"`
	OriginalName  string    `json:"original_filename"`
	ContentType   string    `json:"content_type"`
	Size          int64     `json:"size"`
	CreatedAt     time.Time `json:"created_at"`
	Status        string    `json:"status"`

	// Анализ
	Diagnosis   *int       `json:"diagnosis,omitempty"`
	Confidence  *float32   `json:"confidence,omitempty"`
	GradCAMURL  *string    `json:"gradcam_url,omitempty"`
	CompletedAt *time.Time `json:"completed_at,omitempty"`
}

type UserRefreshResponse struct {
	AccessToken  string `json:"access" binding:"required"`
	RefreshToken string `json:"refresh" binding:"required"`
}

type UserProfile struct {
	ID       int64  `json:"id"`
	Username string `json:"username"`
	Email    string `json:"email"`
	Role     string `json:"role"`
	Status   string `json:"status"`
}
