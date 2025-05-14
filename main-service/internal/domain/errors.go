package domain

import "fmt"

type AppError struct {
	Code    string
	Message string
	Err     error
}

func (e *AppError) Error() string {
	if e.Err != nil {
		return fmt.Sprintf("%s: %v", e.Code, e.Err)
	}
	return e.Code
}

func (e *AppError) Unwrap() error {
	return e.Err
}

func NewAppError(code, message string, err error) *AppError {
	return &AppError{Code: code, Message: message, Err: err}
}

var (
	ErrWrongCredentials = NewAppError(
		"WRONG_CREDENTIALS",
		"Incorrect email or password",
		nil)
	ErrUnauthorized = NewAppError(
		"UNAUTHORIZED",
		"Authorization required or token is invalid",
		nil,
	)
	ErrInternal = func(err error) *AppError {
		return NewAppError(
			"INTERNAL_ERROR",
			"Unexpected server error",
			err,
		)
	}
)

var (
	ErrUserNotFound  = NewAppError("USER_NOT_FOUND", "User not found", nil)
	ErrTokenNotFound = NewAppError("TOKEN_NOT_FOUND", "Token not found", nil)
)

var (
	ErrInviteAlreadyExists = NewAppError("INVITE_EXISTS", "An invite for this email already exists", nil)
	ErrInvalidRole         = NewAppError("INVALID_ROLE", "Invalid user role specified", nil)
)
var (
	ErrEmailAlreadyUsed = NewAppError("INVITE_EXISTS", "Invite already exists for this email", nil)
)
var (
	ErrUserAlreadyExists = NewAppError("USER_EXISTS", "User with this email already exists", nil)
)
var (
	ErrWrongInviteCode   = NewAppError("INVALID_INVITE", "Invalid invite code", nil)
	ErrInviteExpired     = NewAppError("INVITE_EXPIRED", "Invite has expired", nil)
	ErrInviteAlreadyUsed = NewAppError("INVITE_USED", "Invite has already been completed", nil)
)
var (
	ErrResetTokenNotFound = NewAppError("RESET_TOKEN_INVALID", "Reset token is invalid", nil)
	ErrResetTokenExpired  = NewAppError("RESET_TOKEN_EXPIRED", "Reset token has expired", nil)
)
var ErrUserSuspended = NewAppError("USER_SUSPENDED", "User account is suspended", nil)
