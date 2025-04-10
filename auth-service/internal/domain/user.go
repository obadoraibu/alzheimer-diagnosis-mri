package domain

import (
	"database/sql"
)

type User struct {
	Id             int64
	Username       string `json:"username" binding:"required"`
	Email          string `json:"email" binding:"required"`
	PasswordHash   string
	Role           string `json:"role"`
	Status         string `json:"status"`
	InviteToken    sql.NullString
	InviteTokenExp sql.NullTime
}

type CreateUserInvite struct {
	Username string `json:"username" binding:"required"`
	Email    string `json:"email" binding:"required"`
	Role     string `json:"role" binding:"required,oneof=ADMIN DOCTOR PATIENT"`
}

type UserSignUpInput struct {
	Name     string `json:"name" binding:"required"`
	Email    string `json:"email" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type UserSignInInput struct {
	Email       string `json:"email" binding:"required"`
	Password    string `json:"password" binding:"required"`
	Fingerprint string `json:"fingerprint" binding:"required"`
}

type UserSignInResponse struct {
	AccessToken  string `json:"access" binding:"required"`
	RefreshToken string `json:"refresh"binding:"required"`
}

type UserRefreshResponse struct {
	AccessToken  string `json:"access" binding:"required"`
	RefreshToken string `json:"refresh" binding:"required"`
}
