package service

import (
	"database/sql"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/obadoraibu/go-auth/internal/domain"
	"github.com/sirupsen/logrus"
)

func (s *Service) CreateUserInvite(c *gin.Context, r *domain.CreateUserInvite) error {

	code := uuid.New()
	duration, err := time.ParseDuration("24h")
	if err != nil {
		return err
	}

	u := &domain.User{
		Username: r.Username,
		Email:    r.Email,
		Role:     r.Role,
		Status:   "invited",
		InviteToken: sql.NullString{
			String: code.String(),
			Valid:  true,
		},
		InviteTokenExp: sql.NullTime{
			Time:  time.Now().Add(duration),
			Valid: true,
		},
	}

	u, err = s.repo.CreateUserInvite(u)
	if err != nil {
		return err
	}

	logrus.Printf("user %d invite has been created", u.Id)

	err = s.emailSender.SendInvEmail(u.Email, code.String())
	if err != nil {
		return err
	}
	return nil
}
