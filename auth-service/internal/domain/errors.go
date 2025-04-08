package domain

import "errors"

var (
	ErrUserAlreadyExists          = errors.New("user with such email already exists")
	ErrWrongEmailOrPassword       = errors.New("invalid email or password")
	ErrWrongEmailConfirmationCode = errors.New("invalid email confirmation code")
	ErrEmailIsNotConfirmed        = errors.New("email has not been confirmed")
)
