package domain

import (
	"database/sql"
)

const (
	UserStatusActive = "active"
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
	Role     string `json:"role" binding:"required,oneof=admin doctor"`
}

type UserProfile struct {
	ID       int64  `json:"id"`
	Username string `json:"username"`
	Email    string `json:"email"`
	Role     string `json:"role"`
	Status   string `json:"status"`
}

// new
type CreateUserInviteInput struct {
	Username string
	Email    string
	Role     string
}

type UserListFilterInput struct {
	Role   string
	Status string
	Limit  int
	Offset int
}

type UserResponse struct {
	ID       int64  `json:"id"`
	Username string `json:"username"`
	Email    string `json:"email"`
	Role     string `json:"role"`
	Status   string `json:"status"`
}

type UpdateUserInput struct {
	ID       int64
	Username string
	Role     string
	Status   string
}

type UpdateUserOutput struct {
	ID       int64  `json:"id"`
	Username string `json:"username"`
	Email    string `json:"email"`
	Role     string `json:"role"`
	Status   string `json:"status"`
}

type DeleteUserInput struct {
	ID int64
}

type GetUserProfileInput struct {
	UserID int64
}

type UserProfileOutput struct {
	ID       int64  `json:"id"`
	Username string `json:"username"`
	Email    string `json:"email"`
	Role     string `json:"role"`
	Status   string `json:"status"`
}
