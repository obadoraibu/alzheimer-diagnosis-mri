package domain

type User struct {
	Id           string
	Name         string `json:"name" binding:"required"`
	Email        string `json:"email" binding:"required"`
	IsConfirmed  bool
	PasswordHash string
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
