package internal

import (
	"github.com/ductruonghoc/DATN_08_2025_Back-end/config"
	"github.com/golang-jwt/jwt/v4"

	"time"
)

// Claims struct for JWT payload
type UserIDClaims struct {
	UserID int `json:"user_id"`
	jwt.RegisteredClaims
}

func JWTGenerator(userID int) (string, error) {
	expirationTime := time.Now().Add(24 * time.Hour)
	claims := &UserIDClaims{
		UserID: 0,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expirationTime),
		},
	}

	//new token gen
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims);

	jwtKey := config.GetEnv("JWT_KEY", "");

	return token.SignedString([]byte(jwtKey));
}
