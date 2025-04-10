package service

import (
	"github.com/gin-gonic/gin"
	"github.com/obadoraibu/go-auth/internal/domain"
	"github.com/obadoraibu/go-auth/pkg/hash"
)

// new
func (s *Service) SignIn(c *gin.Context, req *domain.UserSignInInput) (*domain.UserSignInResponse, error) {
	u, err := s.repo.FindUserByEmail(req.Email)
	if err != nil {
		return nil, err
	}
	if u.Status != "active" {
		return nil, domain.ErrEmailIsNotConfirmed
	}

	if !hash.CheckPasswordHash(req.Password, u.PasswordHash) {
		return nil, domain.ErrWrongEmailOrPassword
	}

	accessToken, err := s.tokenManager.GenerateJWT(req.Email, u.Role)
	if err != nil {
		return nil, err
	}

	refreshToken := s.tokenManager.GenerateRefresh()

	if err := s.repo.AddToken(req.Fingerprint, refreshToken, req.Email, u.Role); err != nil {
		return nil, err
	}

	response := &domain.UserSignInResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}

	return response, nil
}

// func (s *Service) SignUp(c *gin.Context, inp *domain.UserSignUpInput) error {

// 	hashesPassword, err := hash.HashPassword(inp.Password)
// 	if err != nil {
// 		return err
// 	}

// 	u := &domain.User{
// 		Name:         inp.Name,
// 		Email:        inp.Email,
// 		PasswordHash: hashesPassword,
// 		IsConfirmed:  false,
// 	}

// 	code := uuid.New()
// 	duration, err := time.ParseDuration("5m")
// 	if err != nil {
// 		return err
// 	}

// 	u, err = s.repo.CreateUserAndEmailConfirmation(u, code.String(), time.Now().Add(duration))
// 	if err != nil {
// 		return err
// 	}

// 	logrus.Printf("user %s has been created", u.Id)

// 	err = s.emailSender.SendConfirmationEmail(u.Email, code.String())
// 	if err != nil {
// 		return err
// 	}
// 	return nil
// }

// new
func (s *Service) Refresh(refresh, fingerprint string) (*domain.UserRefreshResponse, error) {
	email, err := s.repo.FindAndDeleteRefreshToken(refresh, fingerprint)
	if err != nil {
		return nil, err
	}

	u, err := s.repo.FindUserByEmail(email)
	if err != nil {
		return nil, err
	}

	access, err := s.tokenManager.GenerateJWT(email, u.Role)
	if err != nil {
		return nil, err
	}

	newRefresh := s.tokenManager.GenerateRefresh()

	if err := s.repo.AddToken(fingerprint, newRefresh, email, u.Role); err != nil {
		return nil, err
	}

	return &domain.UserRefreshResponse{
		RefreshToken: newRefresh,
		AccessToken:  access,
	}, nil
}

func (s *Service) Revoke(refresh, fingerprint string) error {
	_, err := s.repo.FindAndDeleteRefreshToken(refresh, fingerprint)
	if err != nil {
		return err
	}

	return nil
}

func (s *Service) UserInfo(email string) (*domain.User, error) {
	u, err := s.repo.FindUserByEmail(email)
	if err != nil {
		return nil, err
	}

	return u, nil
}

func (s *Service) CompleteInvite(code string, password string) error {

	hashesPassword, err := hash.HashPassword(password)
	if err != nil {
		return err
	}

	err = s.repo.CompleteInvite(code, hashesPassword)
	if err != nil {
		return err
	}

	return nil
}
