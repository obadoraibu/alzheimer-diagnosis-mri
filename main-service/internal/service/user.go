package service

import (
	"database/sql"
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/obadoraibu/go-auth/internal/domain"
	"github.com/sirupsen/logrus"
)

func (s *Service) CreateUserInvite(input *domain.CreateUserInviteInput) error {
	code := uuid.New()
	expiration := time.Now().Add(24 * time.Hour)

	u := &domain.User{
		Username: input.Username,
		Email:    input.Email,
		Role:     input.Role,
		Status:   "invited",
		InviteToken: sql.NullString{
			String: code.String(),
			Valid:  true,
		},
		InviteTokenExp: sql.NullTime{
			Time:  expiration,
			Valid: true,
		},
	}

	created, err := s.repo.CreateUserInvite(u)
	if err != nil {
		if errors.Is(err, domain.ErrUserAlreadyExists) {
			return domain.ErrInviteAlreadyExists
		}
		return domain.ErrInternal(err)
	}

	logrus.WithField("user_id", created.Id).Info("user invite created")

	if err := s.emailSender.SendInvEmail(created.Email, code.String()); err != nil {
		return domain.ErrInternal(err)
	}

	return nil
}

func (s *Service) GetUsersList(input *domain.UserListFilterInput) ([]*domain.UserResponse, error) {
	users, err := s.repo.GetUsersFiltered(input.Role, input.Status, input.Limit, input.Offset)
	if err != nil {
		return nil, domain.ErrInternal(err)
	}

	var resp []*domain.UserResponse
	for _, u := range users {
		resp = append(resp, &domain.UserResponse{
			ID:       u.Id,
			Username: u.Username,
			Email:    u.Email,
			Role:     u.Role,
			Status:   u.Status,
		})
	}

	return resp, nil
}

func (s *Service) UpdateUser(input *domain.UpdateUserInput) error {
	user, err := s.repo.GetUserForUpdate(input.ID)
	if err != nil {
		if errors.Is(err, domain.ErrUserNotFound) {
			return domain.ErrUserNotFound
		}
		return domain.ErrInternal(err)
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

	if err := s.repo.UpdateUserByID(user); err != nil {
		return domain.ErrInternal(err)
	}

	return nil
}

func (s *Service) DeleteUser(input *domain.DeleteUserInput) error {
	user, err := s.repo.GetUserForUpdate(input.ID)
	if err != nil {
		if errors.Is(err, domain.ErrUserNotFound) {
			return domain.ErrUserNotFound
		}
		return domain.ErrInternal(err)
	}

	user.Status = "suspended"

	if err := s.repo.UpdateUserByID(user); err != nil {
		return domain.ErrInternal(err)
	}

	return nil
}

func (s *Service) GetUserProfile(input *domain.GetUserProfileInput) (*domain.UserProfileOutput, error) {
	user, err := s.repo.GetUserByID(input.UserID)
	if err != nil {
		if errors.Is(err, domain.ErrUserNotFound) {
			return nil, domain.ErrUserNotFound
		}
		return nil, domain.ErrInternal(err)
	}

	return &domain.UserProfileOutput{
		ID:       user.Id,
		Username: user.Username,
		Email:    user.Email,
		Role:     user.Role,
		Status:   user.Status,
	}, nil
}
