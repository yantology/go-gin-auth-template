package utils

import (
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type JWTUtil struct {
	accessSecret   string
	refreshSecret  string
	accessTimeout  int
	refreshTimeout int
}

func NewJWTUtil(accessSecret, refreshSecret string, accessTimeout, refreshTimeout int) *JWTUtil {
	return &JWTUtil{
		accessSecret:   accessSecret,
		refreshSecret:  refreshSecret,
		accessTimeout:  accessTimeout,
		refreshTimeout: refreshTimeout,
	}
}

func (j *JWTUtil) GenerateTokens(userID int) (string, string, error) {
	// Generate Access Token
	accessToken := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id": userID,
		"exp":     time.Now().Add(time.Minute * time.Duration(j.accessTimeout)).Unix(),
	})

	accessTokenString, err := accessToken.SignedString([]byte(j.accessSecret))
	if err != nil {
		return "", "", err
	}

	// Generate Refresh Token
	refreshToken := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id": userID,
		"exp":     time.Now().Add(time.Minute * time.Duration(j.refreshTimeout)).Unix(),
	})

	refreshTokenString, err := refreshToken.SignedString([]byte(j.refreshSecret))
	if err != nil {
		return "", "", err
	}

	return accessTokenString, refreshTokenString, nil
}

func (j *JWTUtil) ValidateToken(tokenString string, isRefresh bool) (*jwt.Token, error) {
	secret := j.accessSecret
	if isRefresh {
		secret = j.refreshSecret
	}

	return jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, jwt.ErrSignatureInvalid
		}
		return []byte(secret), nil
	})
}

func (j *JWTUtil) ExtractUserID(token *jwt.Token) (int, error) {
	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return 0, jwt.ErrInvalidKey
	}

	userID := int(claims["user_id"].(float64))
	return userID, nil
}
