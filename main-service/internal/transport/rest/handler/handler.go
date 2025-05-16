package handler

import (
	"context"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/obadoraibu/go-auth/internal/domain"
)

type Service interface {
	SignIn(*domain.UserSignInInput) (*domain.UserSignInOutput, error)
	Refresh(*domain.TokenRefreshInput) (*domain.TokenRefreshOutput, error)
	Revoke(*domain.TokenRevokeInput) error
	CompleteInvite(*domain.CompleteInviteInput) error

	ResetPassword(*domain.ResetPasswordInput) error
	ResetPasswordComplete(*domain.ResetPasswordConfirmInput) error

	CreateUserInvite(*domain.CreateUserInviteInput) error
	GetUsersList(*domain.UserListFilterInput) ([]*domain.UserResponse, error)
	UpdateUser(*domain.UpdateUserInput) error
	DeleteUser(*domain.DeleteUserInput) error

	GetUserProfile(*domain.GetUserProfileInput) (*domain.UserProfileOutput, error)

	//UserInfo(*domain.GetUserProfileInput) (*domain.User, error)

	CreateMRIAnalysis(ctx *gin.Context, input *domain.CreateAnalysisInput) error
	GetScansByFilters(
		ctx context.Context,
		userID int64,
		filter *domain.ScanFilter,
	) ([]*domain.MRIScan, error)
	GetScanByID(
		ctx context.Context,
		userID, scanID int64,
	) (*domain.MRIScanDetail, error)
}

type Handler struct {
	service      Service
	tokenManager TokenManager
}

type TokenManager interface {
	GetSigningKey() string
}

type Dependencies struct {
	Service      Service
	TokenManager TokenManager
}

func NewHandler(deps Dependencies) *Handler {
	return &Handler{
		service:      deps.Service,
		tokenManager: deps.TokenManager,
	}
}

func (h *Handler) InitRoutes() *gin.Engine {
	router := gin.New()

	// Configure CORS middleware
	config := cors.Config{
		AllowOrigins: []string{"http://localhost:3000"},
		AllowMethods: []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders: []string{
			"Authorization",
			"Content-Type",
			"X-Requested-With",
			"Accept",
			"Origin",
			"Access-Control-Allow-Credentials",
		},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
	}

	router.Use(cors.New(config))

	auth := router.Group("/")
	auth.Use(h.AuthMiddleware())
	{
		adminGroup := auth.Group("/admin")
		adminGroup.Use(h.AdminMiddleware())
		{
			adminGroup.POST("/users", h.CreateUserInvite)
			adminGroup.GET("/users", h.ListUsers)
			adminGroup.PUT("/users/:id", h.UpdateUser)
			adminGroup.DELETE("/users/:id", h.DeleteUser)
		}
		auth.POST("/upload", h.CreateMRIAnalysis)
		auth.GET("/scans", h.ListScans)
		auth.GET("/scans/:id", h.GetScanDetail)
		auth.GET("/profile", h.GetProfileInfo)
	}

	router.POST("/reset-password", h.ResetPassword)
	router.POST("/reset-password/:code", h.ResetPasswordConfirm)

	router.POST("/complete-invite/:code", h.CompleteInvite)

	//
	router.POST("/sign-in", h.SignIn)
	router.POST("/refresh", h.Refresh)
	router.POST("/revoke", h.Revoke)

	return router
}
