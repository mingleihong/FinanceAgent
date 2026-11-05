package model

import "github.com/golang-jwt/jwt/v5"

type AccessTokenClaims struct {
	Username string `json:"username"`
	jwt.RegisteredClaims
}

type LoginResult struct {
	AccessToken string `json:"access_token"`
	TokenType   string `json:"token_type"`
}
