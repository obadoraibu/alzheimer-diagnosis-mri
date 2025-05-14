package handler

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/obadoraibu/go-auth/internal/domain"
	"github.com/sirupsen/logrus"
)

func (h *Handler) SignIn(c *gin.Context) {
	type request struct {
		Email       string `json:"email" binding:"required"`
		Password    string `json:"password" binding:"required"`
		Fingerprint string `json:"fingerprint" binding:"required"`
	}

	var req request
	if err := c.ShouldBindJSON(&req); err != nil {
		sendErrorResponse(c, http.StatusBadRequest, "INVALID_JSON", "Invalid request body format")
		return
	}

	input := &domain.UserSignInInput{
		Email:       req.Email,
		Password:    req.Password,
		Fingerprint: req.Fingerprint,
	}

	output, err := h.service.SignIn(input)
	if err != nil {
		var appErr *domain.AppError
		if errors.As(err, &appErr) {
			status := httpStatusFromAppError(appErr)
			sendErrorResponse(c, status, appErr.Code, appErr.Message)
		} else {
			sendErrorResponse(c, http.StatusInternalServerError, "INTERNAL_ERROR", "Unexpected server error")
		}
		return
	}

	cookie := &http.Cookie{
		Name:     "refresh",
		Value:    output.RefreshToken,
		Path:     "/",
		MaxAge:   86400 * 60,
		HttpOnly: true,
		SameSite: http.SameSiteLaxMode,
		// Secure: true,
	}
	http.SetCookie(c.Writer, cookie)

	type response struct {
		AccessToken string `json:"accessToken"`
	}

	c.JSON(http.StatusOK, APIResponse{
		Success: true,
		Data:    response{AccessToken: output.AccessToken},
	})
}

func (h *Handler) Refresh(c *gin.Context) {

	type request struct {
		Fingerprint string `json:"fingerprint" binding:"required"`
	}

	req := &request{}
	if err := c.ShouldBindJSON(&req); err != nil {
		sendErrorResponse(c, http.StatusBadRequest, "INVALID_JSON", "Invalid request body format")
		return
	}

	refresh, err := c.Cookie("refresh")
	if err != nil {
		sendErrorResponse(c, http.StatusUnauthorized, "NO_AUTH_COOKIE", "Authorization cookie is missing")
		return
	}

	input := &domain.TokenRefreshInput{
		Fingerprint: req.Fingerprint,
		Refresh:     refresh,
	}

	output, err := h.service.Refresh(input)
	if err != nil {
		var appErr *domain.AppError
		if errors.As(err, &appErr) {
			status := httpStatusFromAppError(appErr)
			sendErrorResponse(c, status, appErr.Code, appErr.Message)
		} else {
			sendErrorResponse(c, http.StatusInternalServerError, "INTERNAL_ERROR", "Unexpected server error")
		}
		return
	}

	newCookie := &http.Cookie{
		Name:     "refresh",
		Value:    output.RefreshToken,
		Path:     "/",
		MaxAge:   86400 * 60,
		HttpOnly: true,
		SameSite: http.SameSiteLaxMode,
		// Secure: true,
	}

	http.SetCookie(c.Writer, newCookie)

	type response struct {
		AccessToken string `json:"accessToken"`
	}

	c.JSON(http.StatusOK, APIResponse{
		Success: true,
		Data:    response{AccessToken: output.AccessToken},
	})
}

func (h *Handler) Revoke(c *gin.Context) {
	type request struct {
		Fingerprint string `json:"fingerprint" binding:"required"`
	}

	var req request
	if err := c.ShouldBindJSON(&req); err != nil {
		sendErrorResponse(c, http.StatusBadRequest, "INVALID_JSON", "Invalid request body format")
		return
	}

	refresh, err := c.Cookie("refresh")
	if err != nil {
		sendErrorResponse(c, http.StatusUnauthorized, "NO_AUTH_COOKIE", "Authorization cookie is missing")
		return
	}

	input := &domain.TokenRevokeInput{
		Fingerprint: req.Fingerprint,
		Refresh:     refresh,
	}

	err = h.service.Revoke(input)
	if err != nil {
		var appErr *domain.AppError
		if errors.As(err, &appErr) {
			status := httpStatusFromAppError(appErr)
			sendErrorResponse(c, status, appErr.Code, appErr.Message)
		} else {
			sendErrorResponse(c, http.StatusInternalServerError, "INTERNAL_ERROR", "Unexpected server error")
		}
		return
	}

	http.SetCookie(c.Writer, &http.Cookie{
		Name:     "refresh",
		Value:    "",
		Path:     "/",
		MaxAge:   -1,
		HttpOnly: true,
	})

	c.JSON(http.StatusOK, APIResponse{
		Success: true,
	})
}

func (h *Handler) CompleteInvite(c *gin.Context) {
	type request struct {
		Password string `json:"password" binding:"required"`
	}

	var req request
	if err := c.ShouldBindJSON(&req); err != nil {
		sendErrorResponse(c, http.StatusBadRequest, "INVALID_JSON", "Invalid request body format")
		return
	}

	code := c.Param("code")
	if code == "" {
		sendErrorResponse(c, http.StatusBadRequest, "MISSING_INVITE_CODE", "Invite code is required")
		return
	}

	input := &domain.CompleteInviteInput{
		Code:     code,
		Password: req.Password,
	}

	err := h.service.CompleteInvite(input)
	if err != nil {
		var appErr *domain.AppError
		if errors.As(err, &appErr) {
			status := httpStatusFromAppError(appErr)
			sendErrorResponse(c, status, appErr.Code, appErr.Message)
		} else {
			sendErrorResponse(c, http.StatusInternalServerError, "INTERNAL_ERROR", "Unexpected server error")
		}
		return
	}

	c.JSON(http.StatusOK, APIResponse{
		Success: true,
	})
}

func (h *Handler) ResetPassword(c *gin.Context) {
	type request struct {
		Email string `json:"email" binding:"required,email"`
	}

	var req request
	if err := c.ShouldBindJSON(&req); err != nil {
		sendErrorResponse(c, http.StatusBadRequest, "INVALID_JSON", "Invalid request body")
		return
	}

	input := &domain.ResetPasswordInput{
		Email: req.Email,
	}

	err := h.service.ResetPassword(input)
	if err != nil {
		var appErr *domain.AppError
		if errors.As(err, &appErr) {
			sendErrorResponse(c, httpStatusFromAppError(appErr), appErr.Code, appErr.Message)
		} else {
			sendErrorResponse(c, http.StatusInternalServerError, "INTERNAL_ERROR", "Unexpected server error")
		}
		return
	}

	c.JSON(http.StatusOK, APIResponse{
		Success: true,
	})
}

func (h *Handler) ResetPasswordConfirm(c *gin.Context) {
	type request struct {
		Password string `json:"password" binding:"required"`
	}

	var req request
	if err := c.ShouldBindJSON(&req); err != nil {
		sendErrorResponse(c, http.StatusBadRequest, "INVALID_JSON", "Invalid request body format")
		return
	}

	token := c.Param("code")
	if token == "" {
		sendErrorResponse(c, http.StatusBadRequest, "MISSING_TOKEN", "Reset token is required")
		return
	}

	input := &domain.ResetPasswordConfirmInput{
		Token:    token,
		Password: req.Password,
	}

	if err := h.service.ResetPasswordComplete(input); err != nil {
		var appErr *domain.AppError
		if errors.As(err, &appErr) {
			sendErrorResponse(c, httpStatusFromAppError(appErr), appErr.Code, appErr.Message)
		} else {
			sendErrorResponse(c, http.StatusInternalServerError, "INTERNAL_ERROR", "Unexpected server error")
		}
		return
	}

	logrus.Println("password changed")
	c.Redirect(http.StatusFound, "http://localhost:3000/sign-in")
}
