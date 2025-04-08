package handler

import (
	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"net/http"
	"strings"
)

func (h *Handler) AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {

		tokenString := c.GetHeader("Authorization")

		if tokenString == "" {
			sendErrorResponse(c, http.StatusUnauthorized, "cannot read authorization token")
			c.Abort()
			return
		}

		// Split the "Authorization" header to get the token part after "Bearer "
		parts := strings.Split(tokenString, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			sendErrorResponse(c, http.StatusUnauthorized, "invalid authorization header format")
			c.Abort()
			return
		}

		token, err := jwt.Parse(parts[1], func(token *jwt.Token) (interface{}, error) {
			return []byte(h.tokenManager.GetSigningKey()), nil
		})

		if err != nil || !token.Valid {
			sendErrorResponse(c, http.StatusUnauthorized, err.Error())
			c.Abort()
			return
		}

		c.Set("AccessToken", token)
		c.Next()
	}
}

func CORSMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {

		c.Header("Access-Control-Allow-Origin", "*")
		c.Header("Access-Control-Allow-Credentials", "true")
		c.Header("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
		c.Header("Access-Control-Allow-Methods", "POST,HEAD,PATCH, OPTIONS, GET, PUT")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	}
}
