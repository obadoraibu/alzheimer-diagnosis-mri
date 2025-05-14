package handler

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"github.com/obadoraibu/go-auth/internal/domain"
)

func (h *Handler) CreateUserInvite(c *gin.Context) {
	type request struct {
		Username string `json:"username" binding:"required"`
		Email    string `json:"email" binding:"required,email"`
		Role     string `json:"role" binding:"required"`
	}

	var req request
	if err := c.ShouldBindJSON(&req); err != nil {
		sendErrorResponse(c, http.StatusBadRequest, "INVALID_JSON", "Invalid request body format")
		return
	}

	input := &domain.CreateUserInviteInput{
		Username: req.Username,
		Email:    req.Email,
		Role:     req.Role,
	}

	err := h.service.CreateUserInvite(input)
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

func (h *Handler) ListUsers(c *gin.Context) {
	limitStr := c.DefaultQuery("limit", "20")
	offsetStr := c.DefaultQuery("offset", "0")

	limit, err := strconv.Atoi(limitStr)
	if err != nil || limit <= 0 {
		sendErrorResponse(c, http.StatusBadRequest, "INVALID_LIMIT", "Limit must be a positive integer")
		return
	}

	offset, err := strconv.Atoi(offsetStr)
	if err != nil || offset < 0 {
		sendErrorResponse(c, http.StatusBadRequest, "INVALID_OFFSET", "Offset must be a non-negative integer")
		return
	}

	input := &domain.UserListFilterInput{
		Role:   c.Query("role"),
		Status: c.Query("status"),
		Limit:  limit,
		Offset: offset,
	}

	users, err := h.service.GetUsersList(input)
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
		Data:    users,
		Meta: &Meta{
			Limit:  limit,
			Offset: offset,
			Count:  len(users),
		},
	})
}

func (h *Handler) UpdateUser(c *gin.Context) {
	type request struct {
		Username string `json:"username,omitempty"`
		Role     string `json:"role,omitempty"`
		Status   string `json:"status,omitempty"`
	}

	idParam := c.Param("id")
	userID, err := strconv.ParseInt(idParam, 10, 64)
	if err != nil {
		sendErrorResponse(c, http.StatusBadRequest, "INVALID_USER_ID", "User ID must be a valid integer")
		return
	}

	var req request
	if err := c.ShouldBindJSON(&req); err != nil {
		sendErrorResponse(c, http.StatusBadRequest, "INVALID_JSON", "Invalid request body")
		return
	}

	input := &domain.UpdateUserInput{
		ID:       userID,
		Username: req.Username,
		Role:     req.Role,
		Status:   req.Status,
	}

	err = h.service.UpdateUser(input)
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

func (h *Handler) DeleteUser(c *gin.Context) {
	idParam := c.Param("id")
	userID, err := strconv.ParseInt(idParam, 10, 64)
	if err != nil {
		sendErrorResponse(c, http.StatusBadRequest, "INVALID_USER_ID", "User ID must be a valid integer")
		return
	}

	input := &domain.DeleteUserInput{
		ID: userID,
	}

	err = h.service.DeleteUser(input)
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

func (h *Handler) GetProfileInfo(c *gin.Context) {
	token := c.MustGet("AccessToken").(*jwt.Token)
	claims := token.Claims.(jwt.MapClaims)

	userID, ok := claims["user_id"].(float64)
	if !ok {
		sendErrorResponse(c, http.StatusUnauthorized, "INVALID_TOKEN", "Invalid access token payload")
		return
	}

	input := &domain.GetUserProfileInput{
		UserID: int64(userID),
	}

	profile, err := h.service.GetUserProfile(input)
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
		Data:    profile,
	})
}
