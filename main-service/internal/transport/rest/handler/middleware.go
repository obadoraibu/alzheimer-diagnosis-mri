package handler

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
)

func (h *Handler) AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			sendErrorResponse(c, http.StatusUnauthorized, "NO_AUTH_HEADER", "Authorization header is missing")
			c.Abort()
			return
		}

		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			sendErrorResponse(c, http.StatusUnauthorized, "INVALID_AUTH_HEADER", "Authorization header must be in Bearer format")
			c.Abort()
			return
		}

		tokenStr := parts[1]
		token, err := jwt.Parse(tokenStr, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
			}
			return []byte(h.tokenManager.GetSigningKey()), nil
		})

		if err != nil || !token.Valid {
			sendErrorResponse(c, http.StatusUnauthorized, "INVALID_TOKEN", "Invalid or expired token")
			c.Abort()
			return
		}

		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			sendErrorResponse(c, http.StatusUnauthorized, "INVALID_CLAIMS", "Invalid token claims")
			c.Abort()
			return
		}

		// Извлечение user_id
		rawUserID, ok := claims["user_id"]
		if !ok {
			sendErrorResponse(c, http.StatusUnauthorized, "MISSING_USER_ID", "User ID missing in token")
			c.Abort()
			return
		}

		var userID int64
		switch v := rawUserID.(type) {
		case float64:
			userID = int64(v)
		case string:
			parsedID, err := strconv.ParseInt(v, 10, 64)
			if err != nil {
				sendErrorResponse(c, http.StatusUnauthorized, "INVALID_USER_ID", "User ID must be numeric")
				c.Abort()
				return
			}
			userID = parsedID
		default:
			sendErrorResponse(c, http.StatusUnauthorized, "INVALID_USER_ID_TYPE", "Unexpected user ID type in token")
			c.Abort()
			return
		}

		// Извлечение роли
		role, ok := claims["role"].(string)
		if !ok {
			sendErrorResponse(c, http.StatusUnauthorized, "MISSING_ROLE", "Role missing in token")
			c.Abort()
			return
		}

		// Установка в контекст
		c.Set("UserID", userID)
		c.Set("Role", role)
		c.Set("AccessToken", token)

		c.Next()
	}
}

// func CORSMiddleware() gin.HandlerFunc {
// 	return func(c *gin.Context) {

// 		c.Header("Access-Control-Allow-Origin", "*")
// 		c.Header("Access-Control-Allow-Credentials", "true")
// 		c.Header("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
// 		c.Header("Access-Control-Allow-Methods", "POST,HEAD,PATCH, OPTIONS, GET, PUT")

// 		if c.Request.Method == "OPTIONS" {
// 			c.AbortWithStatus(204)
// 			return
// 		}

// 		c.Next()
// 	}
// }

func (h *Handler) AdminMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		roleInterface, exists := c.Get("Role")
		if !exists {
			sendErrorResponse(c, http.StatusUnauthorized, "NO_ROLE_IN_CONTEXT", "Missing role in request context")
			c.Abort()
			return
		}

		role, ok := roleInterface.(string)
		if !ok || role != "admin" {
			sendErrorResponse(c, http.StatusForbidden, "ADMIN_REQUIRED", "Admin access required")
			c.Abort()
			return
		}

		c.Next()
	}
}
