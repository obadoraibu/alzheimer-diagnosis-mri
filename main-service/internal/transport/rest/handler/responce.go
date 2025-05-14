package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/obadoraibu/go-auth/internal/domain"
	"github.com/sirupsen/logrus"
)

type errorResponse struct {
	Message string `json:"message"`
}

type APIError struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

type Meta struct {
	Limit  int `json:"limit"`
	Offset int `json:"offset"`
	Count  int `json:"count"`
}

type APIResponse struct {
	Success bool        `json:"success"`
	Error   *APIError   `json:"error,omitempty"`
	Data    interface{} `json:"data,omitempty"`
	Meta    *Meta       `json:"meta,omitempty"`
}

func sendErrorResponse(c *gin.Context, statusCode int, code string, message string) {
	logrus.WithFields(logrus.Fields{
		"code":    code,
		"message": message,
	}).Error("API error")

	c.AbortWithStatusJSON(statusCode, APIResponse{
		Success: false,
		Error: &APIError{
			Code:    code,
			Message: message,
		},
	})
}

var errorStatusMap = map[string]int{
	"WRONG_CREDENTIALS": http.StatusUnauthorized,
	"INVALID_JSON":      http.StatusBadRequest,
	"INTERNAL_ERROR":    http.StatusInternalServerError,
	"USER_SUSPENDED":    http.StatusForbidden,
}

func httpStatusFromAppError(err *domain.AppError) int {
	if status, ok := errorStatusMap[err.Code]; ok {
		return status
	}
	return http.StatusInternalServerError
}
