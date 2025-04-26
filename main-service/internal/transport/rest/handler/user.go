package handler

import (
	"net/http"
	"strconv"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"github.com/obadoraibu/go-auth/internal/domain"
)

func (h *Handler) CreateUserInvite(c *gin.Context) {
	r := &domain.CreateUserInvite{}
	if err := c.ShouldBindJSON(&r); err != nil {
		sendErrorResponse(c, http.StatusBadRequest, "invalid request body")
		return
	}

	err := h.service.CreateUserInvite(r)
	if err != nil {
		sendErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.Status(http.StatusOK)
}

func (h *Handler) ListUsers(c *gin.Context) {
	// Получаем query-параметры
	role := c.Query("role")
	status := c.Query("status")

	// Пагинация с дефолтами
	limitStr := c.DefaultQuery("limit", "20")
	offsetStr := c.DefaultQuery("offset", "0")

	limit, err := strconv.Atoi(limitStr)
	if err != nil || limit <= 0 {
		sendErrorResponse(c, http.StatusBadRequest, "invalid limit")
		return
	}

	offset, err := strconv.Atoi(offsetStr)
	if err != nil || offset < 0 {
		sendErrorResponse(c, http.StatusBadRequest, "invalid offset")
		return
	}

	users, err := h.service.GetUsersList(role, status, limit, offset)
	if err != nil {
		sendErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	var response []*domain.UserResponse
	for _, u := range users {
		response = append(response, &domain.UserResponse{
			ID:       u.Id,
			Username: u.Username,
			Email:    u.Email,
			Role:     u.Role,
			Status:   u.Status,
		})
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "success",
		"data":    response,
		"meta": gin.H{
			"limit":  limit,
			"offset": offset,
			"count":  len(response),
		},
	})
}

func (h *Handler) UpdateUser(c *gin.Context) {
	idParam := c.Param("id")
	userID, err := strconv.ParseInt(idParam, 10, 64)
	if err != nil {
		sendErrorResponse(c, http.StatusBadRequest, "invalid user ID")
		return
	}

	r := &domain.UpdateUserInput{}
	if err := c.ShouldBindJSON(&r); err != nil {
		sendErrorResponse(c, http.StatusBadRequest, "invalid JSON body")
		return
	}

	err = h.service.UpdateUser(userID, r)
	if err != nil {
		sendErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "user updated"})
}

func (h *Handler) DeleteUser(c *gin.Context) {
	idParam := c.Param("id")
	userID, err := strconv.ParseInt(idParam, 10, 64)
	if err != nil {
		sendErrorResponse(c, http.StatusBadRequest, "invalid user ID")
		return
	}

	err = h.service.DeleteUser(userID)
	if err != nil {
		sendErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "user deleted"})
}

func (h *Handler) GetProfileInfo(c *gin.Context) {
	token := c.MustGet("AccessToken").(*jwt.Token)
	claims := token.Claims.(jwt.MapClaims)
	userID := int64(claims["user_id"].(float64))

	profile, err := h.service.GetUserProfile(userID)
	if err != nil {
		sendErrorResponse(c, http.StatusInternalServerError, "failed to get user profile")
		return
	}

	c.JSON(http.StatusOK, profile)
}
