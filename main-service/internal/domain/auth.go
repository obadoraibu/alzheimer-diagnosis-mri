package domain

type UserSignInInput struct {
	Email       string `json:"email" binding:"required"`
	Password    string `json:"password" binding:"required"`
	Fingerprint string `json:"fingerprint" binding:"required"`
}

type UserSignInOutput struct {
	AccessToken  string `json:"access" binding:"required"`
	RefreshToken string `json:"refresh" binding:"required"`
}

// token (refresh/revoke)
type TokenData struct {
	Id   string `json:"user_id"`
	Role string `json:"role"`
}

type TokenRefreshInput struct {
	Fingerprint string `json:"fingerprint" binding:"required"`
	Refresh     string
}

type TokenRefreshOutput struct {
	AccessToken  string `json:"access" binding:"required"`
	RefreshToken string `json:"refresh" binding:"required"`
}

type TokenRevokeInput struct {
	Refresh     string
	Fingerprint string
}

// complete invite
type CompleteInviteInput struct {
	Code     string
	Password string
}

// reset passsword
type ResetPasswordInput struct {
	Email string
}

type ResetPasswordConfirmInput struct {
	Token    string
	Password string
}
