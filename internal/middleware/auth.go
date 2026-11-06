package middleware

import (
	"Finance/internal/config"
	"Finance/internal/model"
	"Finance/internal/response"
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

// JWTAuth validates access tokens before a handler can start its response.
func JWTAuth(cfg *config.JWTConfig) gin.HandlerFunc {
	secret := []byte(cfg.Secret)
	return func(c *gin.Context) {
		headers := c.Request.Header.Values("Authorization")
		if len(headers) != 1 {
			unauthorized(c)
			return
		}
		parts := strings.Fields(headers[0])
		if len(parts) != 2 || !strings.EqualFold(parts[0], "Bearer") {
			unauthorized(c)
			return
		}

		claims := &model.AccessTokenClaims{}
		_, err := jwt.ParseWithClaims(parts[1], claims, func(_ *jwt.Token) (any, error) {
			return secret, nil
		}, jwt.WithValidMethods([]string{jwt.SigningMethodHS256.Alg()}), jwt.WithExpirationRequired())
		if err != nil {
			unauthorized(c)
			return
		}
		userID, err := strconv.ParseInt(claims.Subject, 10, 64)
		if err != nil {
			unauthorized(c)
			return
		}

		c.Set("user_id", userID)
		c.Set("username", claims.Username)
		c.Next()
	}
}

func unauthorized(c *gin.Context) {
	c.Header("WWW-Authenticate", "Bearer")
	c.Header("Cache-Control", "no-store")
	response.Error(c, http.StatusUnauthorized, "invalid or expired access token")
}
