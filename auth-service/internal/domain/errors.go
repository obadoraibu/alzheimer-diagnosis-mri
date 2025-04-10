package domain

import "errors"

var (
	ErrUserAlreadyExists          = errors.New("user with this email already exists")
	ErrWrongEmailOrPassword       = errors.New("invalid email or password")
	ErrWrongEmailConfirmationCode = errors.New("invalid email confirmation code")
	ErrEmailIsNotConfirmed        = errors.New("email has not been confirmed")

	ErrWrongInviteCode            = errors.New("invalid or unknown invite code")
	ErrInviteAlreadyUsed          = errors.New("invite has already been used or account is active")
	ErrInviteExpired              = errors.New("invite link has expired")
)

