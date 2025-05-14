package service

import (
	"context"
	"mime/multipart"
	"net/url"
	"time"

	"github.com/obadoraibu/go-auth/internal/domain"
)

type Service struct {
	repo         Repository
	tokenManager TokenManager
	emailSender  EmailSender
}

type Repository interface {
	FindUserByEmail(email string) (*domain.User, error)
	AddToken(fingerprint, refresh string, user_id int64, role string) error
	DeleteToken(u *domain.User) error
	FindAndDeleteRefreshToken(refresh, fingerprint string) (string, error)
	Close() error
	CreateUserInvite(u *domain.User) (*domain.User, error)
	CompleteInvite(code, passwordHash string) error
	GetUsersFiltered(role, status string, limit, offset int) ([]*domain.User, error)

	GetUserByID(userID int64) (*domain.User, error)
	GetUserForUpdate(userID int64) (*domain.User, error)
	UpdateUserByID(user *domain.User) error

	UploadScanToMinIO(ctx context.Context, objectName string, file multipart.File, size int64, contentType string) error
	SaveScanMetadata(
		userID int64,
		objectName string,
		originalFilename string,
		contentType string,
		size int64,
		patientName string,
		patientGender string,
		patientAge int,
		scanDate time.Time,
	) (int64, error)
	EnqueueScanTask(userID int64, objectName string) error
	GetScansByFilters(userID int64, filter *domain.ScanFilter) ([]*domain.MRIScan, error)
	GetScanDetail(userID, scanID int64) (*domain.MRIScanDetail, error)

	PresignedGetObject(objectName string) (*url.URL, error)

	SaveResetToken(userID int64, resetToken string, expiresAt time.Time) error
	FindUserByResetToken(token string) (*domain.User, error)

	UpdateUserPassword(userID int64, hash string) error
}

type TokenManager interface {
	GenerateJWT(user_id int64, role string) (string, error)
	GenerateRefresh() string
	GenerateResetToken() string
}

type EmailSender interface {
	SendInvEmail(to, code string) error
	SendPasswordResetEmail(to, code string) error
}

type Dependencies struct {
	Repo         Repository
	TokenManager TokenManager
	EmailService EmailSender
}

func NewService(deps Dependencies) *Service {
	return &Service{
		repo:         deps.Repo,
		tokenManager: deps.TokenManager,
		emailSender:  deps.EmailService,
	}
}
