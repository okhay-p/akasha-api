package jwt

import (
	"akasha-api/pkg/config"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type JwtClaims struct {
	jwt.RegisteredClaims
}

var jwtSecret []byte

func SetSecret(cfg config.Config) {
	jwtSecret = []byte(cfg.JwtSecret)
}

func CreateToken(uuid string) (string, error) {
	claims := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub": uuid,
		"iss": "akashalearn",
		"exp": time.Now().Add(24 * time.Hour).Unix(),
		"iat": time.Now().Unix(),
	})

	tokenString, err := claims.SignedString(jwtSecret)
	if err != nil {
		return "", err
	}

	return tokenString, nil

}

func VerifyToken(tokenString string) (*JwtClaims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &JwtClaims{}, func(token *jwt.Token) (any, error) {
		// Check the signing method
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return jwtSecret, nil
	})

	if err != nil {
		return nil, fmt.Errorf("failed to parse token: %w", err)
	}

	if claims, ok := token.Claims.(*JwtClaims); ok && token.Valid {
		return claims, nil
	}

	return nil, fmt.Errorf("invalid token")
}
