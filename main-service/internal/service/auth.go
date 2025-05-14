package service

import (
	"encoding/json"
	"errors"
	"strconv"
	"time"

	"github.com/obadoraibu/go-auth/internal/domain"
	"github.com/obadoraibu/go-auth/pkg/hash"
)

func (s *Service) SignIn(req *domain.UserSignInInput) (*domain.UserSignInOutput, error) {
	u, err := s.repo.FindUserByEmail(req.Email)
	if err != nil {
		if errors.Is(err, domain.ErrUserNotFound) {
			return nil, domain.ErrWrongCredentials
		}
		return nil, domain.ErrInternal(err)
	}

	if err := s.checkUserActive(u.Id); err != nil {
		return nil, err
	}

	if !hash.CheckPasswordHash(req.Password, u.PasswordHash) {
		return nil, domain.ErrWrongCredentials
	}

	accessToken, err := s.tokenManager.GenerateJWT(u.Id, u.Role)
	if err != nil {
		return nil, domain.ErrInternal(err)
	}

	refreshToken := s.tokenManager.GenerateRefresh()

	if err := s.repo.AddToken(req.Fingerprint, refreshToken, u.Id, u.Role); err != nil {
		return nil, domain.ErrInternal(err)
	}

	return &domain.UserSignInOutput{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}, nil
}

func (s *Service) Refresh(req *domain.TokenRefreshInput) (*domain.TokenRefreshOutput, error) {
	raw, err := s.repo.FindAndDeleteRefreshToken(req.Refresh, req.Fingerprint)
	if err != nil {
		if errors.Is(err, domain.ErrTokenNotFound) {
			return nil, domain.ErrUnauthorized
		}
		return nil, domain.ErrInternal(err)
	}

	var tokenData domain.TokenData
	if err := json.Unmarshal([]byte(raw), &tokenData); err != nil {
		return nil, domain.NewAppError("INVALID_REFRESH_DATA", "Corrupted refresh token data", err)
	}

	userID, err := strconv.ParseInt(tokenData.Id, 10, 64)
	if err != nil {
		return nil, domain.NewAppError("INVALID_USER_ID", "Invalid user ID in token", err)
	}

	u, err := s.repo.GetUserByID(userID)
	if err != nil {
		if errors.Is(err, domain.ErrUserNotFound) {
			return nil, domain.ErrWrongCredentials
		}
		return nil, domain.ErrInternal(err)
	}

	if err := s.checkUserActive(userID); err != nil {
		return nil, err
	}

	access, err := s.tokenManager.GenerateJWT(u.Id, u.Role)
	if err != nil {
		return nil, domain.ErrInternal(err)
	}

	newRefresh := s.tokenManager.GenerateRefresh()

	if err := s.repo.AddToken(req.Fingerprint, newRefresh, u.Id, u.Role); err != nil {
		return nil, domain.ErrInternal(err)
	}

	return &domain.TokenRefreshOutput{
		RefreshToken: newRefresh,
		AccessToken:  access,
	}, nil
}

func (s *Service) Revoke(req *domain.TokenRevokeInput) error {
	_, err := s.repo.FindAndDeleteRefreshToken(req.Refresh, req.Fingerprint)
	if err != nil {
		if errors.Is(err, domain.ErrTokenNotFound) {
			return domain.ErrUnauthorized
		}
		return domain.ErrInternal(err)
	}

	return nil
}

func (s *Service) CompleteInvite(input *domain.CompleteInviteInput) error {
	hashedPassword, err := hash.HashPassword(input.Password)
	if err != nil {
		return domain.ErrInternal(err)
	}

	err = s.repo.CompleteInvite(input.Code, hashedPassword)
	if err != nil {
		if errors.Is(err, domain.ErrWrongInviteCode) {
			return domain.ErrUnauthorized
		}
		if errors.Is(err, domain.ErrInviteExpired) || errors.Is(err, domain.ErrInviteAlreadyUsed) {
			return err
		}
		return domain.ErrInternal(err)
	}

	return nil
}

func (s *Service) ResetPassword(input *domain.ResetPasswordInput) error {
	user, err := s.repo.FindUserByEmail(input.Email)
	if err != nil {
		if errors.Is(err, domain.ErrUserNotFound) {
			return domain.ErrUserNotFound
		}
		return domain.ErrInternal(err)
	}
	if err := s.checkUserActive(user.Id); err != nil {
		return err
	}

	resetToken := s.tokenManager.GenerateResetToken()
	expiresAt := time.Now().Add(1 * time.Hour)

	err = s.repo.SaveResetToken(user.Id, resetToken, expiresAt)
	if err != nil {
		return domain.ErrInternal(err)
	}

	if err := s.emailSender.SendPasswordResetEmail(user.Email, resetToken); err != nil {
		return domain.ErrInternal(err)
	}

	return nil
}

func (s *Service) ResetPasswordComplete(input *domain.ResetPasswordConfirmInput) error {
	user, err := s.repo.FindUserByResetToken(input.Token)
	if err != nil {
		if errors.Is(err, domain.ErrResetTokenNotFound) {
			return domain.ErrResetTokenNotFound
		}
		return domain.ErrInternal(err)
	}

	if !user.InviteTokenExp.Valid || user.InviteTokenExp.Time.Before(time.Now()) {
		return domain.ErrResetTokenExpired
	}

	hashedPassword, err := hash.HashPassword(input.Password)
	if err != nil {
		return domain.ErrInternal(err)
	}

	if err := s.repo.UpdateUserPassword(user.Id, hashedPassword); err != nil {
		return domain.ErrInternal(err)
	}

	return nil
}

func (s *Service) checkUserActive(userID int64) error {
	user, err := s.repo.GetUserByID(userID)
	if err != nil {
		if errors.Is(err, domain.ErrUserNotFound) {
			return domain.ErrUserNotFound
		}
		return domain.ErrInternal(err)
	}
	if user.Status == "suspended" {
		return domain.ErrUserSuspended
	}
	return nil
}
