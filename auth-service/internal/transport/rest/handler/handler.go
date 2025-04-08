package handler

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/obadoraibu/go-auth/internal/domain"
)

type Service interface {
	SignIn(c *gin.Context, r *domain.UserSignInInput) (*domain.UserSignInResponse, error)
	SignUp(c *gin.Context, r *domain.UserSignUpInput) error
	Refresh(refresh, fingerprint string) (*domain.UserRefreshResponse, error)
	Revoke(refresh, fingerprint string) error
	UserInfo(email string) (*domain.User, error)
	ConfirmEmail(code string) error
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

	router.POST("/sign-up", h.SignUp)
	router.GET("/resource", h.AuthMiddleware(), h.UserInfo)
	router.POST("/sign-in", h.SignIn)
	router.POST("/refresh", h.Refresh)
	router.POST("/revoke", h.Revoke)
	router.GET("/email-confirm/:code", h.ConfirmEmail)

	return router
}
