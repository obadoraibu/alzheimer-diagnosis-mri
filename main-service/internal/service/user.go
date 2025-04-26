package service

import (
	"database/sql"
	"time"

	"github.com/google/uuid"
	"github.com/obadoraibu/go-auth/internal/domain"
	"github.com/sirupsen/logrus"
)

func (s *Service) CreateUserInvite(r *domain.CreateUserInvite) error {

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

func (s *Service) GetUsersList(role, status string, limit, offset int) ([]*domain.User, error) {
	return s.repo.GetUsersFiltered(role, status, limit, offset)
}

func (s *Service) UpdateUser(userID int64, input *domain.UpdateUserInput) error {

	user, err := s.repo.GetUserByID(userID)
	if err != nil {
		return err
	}

	if input.Username != "" {
		user.Username = input.Username
	}
	if input.Role != "" {
		user.Role = input.Role
	}
	if input.Status != "" {
		user.Status = input.Status
	}

	err = s.repo.UpdateUserByID(user)
	if err != nil {
		return err
	}

	return nil
}

func (s *Service) DeleteUser(userID int64) error {
	user, err := s.repo.GetUserByID(userID)
	if err != nil {
		return err
	}

	user.Status = "suspended"
	return s.repo.UpdateUserByID(user)
}

func (s *Service) GetUserProfile(userID int64) (*domain.UserProfile, error) {
	user, err := s.repo.GetUserByID(userID)
	if err != nil {
		return nil, err
	}

	return &domain.UserProfile{
		ID:       user.Id,
		Username: user.Username,
		Email:    user.Email,
		Role:     user.Role,
		Status:   user.Status,
	}, nil
}