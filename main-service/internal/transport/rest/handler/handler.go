package handler

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/obadoraibu/go-auth/internal/domain"
)

type Service interface {
	SignIn(c *gin.Context, r *domain.UserSignInInput) (*domain.UserSignInResponse, error)
	//SignUp(c *gin.Context, r *domain.UserSignUpInput) error
	Refresh(refresh, fingerprint string) (*domain.UserRefreshResponse, error)
	Revoke(refresh, fingerprint string) error
	UserInfo(email string) (*domain.User, error)
	CompleteInvite(code string, password string) error
	CreateUserInvite(c *gin.Context, r *domain.CreateUserInvite) error
	GetUsersList(c *gin.Context, role, status string, limit, offset int) ([]*domain.User, error)

	UpdateUser(userID int64, input *domain.UpdateUserInput) error
	DeleteUser(userID int64) error

	CreateMRIAnalysis(c *gin.Context, input *domain.CreateAnalysisInput) error
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
		AllowOrigins:     []string{"http://localhost:3000", "http://localhost:3001"}, // Allow frontend origin
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Authorization", "Content-Type"},
		AllowCredentials: true,
	}

	router.Use(cors.New(config))

	auth := router.Group("/")
	auth.Use(h.AuthMiddleware())
	{
		adminGroup := auth.Group("/admin")
		adminGroup.Use(h.AdminMiddleware())
		{
			adminGroup.POST("/users", h.CreateUserInvite) // Пригласить нового пользователя (по умолчанию DOCTOR)
			adminGroup.GET("/users", h.ListUsers)
			//adminGroup.GET("/users/:id", h.GetUserByID)         // Детальная инфа
			adminGroup.PUT("/users/:id", h.UpdateUser)    // Изменить (роль, статус, ФИО и т.д.)
			adminGroup.DELETE("/users/:id", h.DeleteUser) // Удалить/заблокировать и т.п.
		}
		auth.POST("/upload", h.CreateMRIAnalysis)
	}

	//router.POST("/complete-invite", h.CompleteInvite)

	// router.POST("/sign-up", h.SignUp)

	router.POST("/complete-invite/:code", h.CompleteInvite)

	router.POST("/sign-in", h.SignIn)
	router.POST("/refresh", h.Refresh)
	router.POST("/revoke", h.Revoke)
	//router.GET("/email-confirm/:code", h.ConfirmEmail)

	return router
}
