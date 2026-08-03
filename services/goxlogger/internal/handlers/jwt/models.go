package jwt

import "github.com/golang-jwt/jwt/v5"

type User struct {
	ID           string
	Email        string
	PasswordHash []byte
}

type contextKey string

const userContextKey contextKey = "user"

type Claims struct {
	jwt.RegisteredClaims
	Email string `json:"email"`
}

type TokenPair struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}

type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}
