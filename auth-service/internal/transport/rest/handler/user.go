package handler

import (
	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"net/http"
)

func (h *Handler) UserInfo(c *gin.Context) {

	tokenInterface, exists := c.Get("AccessToken")
	if !exists {
		sendErrorResponse(c, http.StatusUnauthorized, "token not found in context")
		return
	}

	token, ok := tokenInterface.(*jwt.Token)
	if !ok {
		sendErrorResponse(c, http.StatusUnauthorized, "invalid token type")
		return
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		sendErrorResponse(c, http.StatusUnauthorized, "unable to extract claims")
		return
	}

	email, ok := claims["email"].(string)
	if !ok {
		sendErrorResponse(c, http.StatusUnauthorized, "email claim is not a string")
		return
	}

	u, err := h.service.UserInfo(email)
	if err != nil {
		sendErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	c.JSON(http.StatusOK, u)
}
