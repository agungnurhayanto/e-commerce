package auth

import (
	"time"

	"github.com/golang-jwt/jwt/v5"

	"github.com/google/uuid"
)

func GenerateToken(userID uuid.UUID,
	secret string,
	expireHours int,
) (string, error) {
	claims := jwt.MapClaims{
		"user_id": userID.String(),
		"exp":     time.Now().Add(time.Duration(expireHours) * time.Hour).Unix(),
		"iat":     time.Now().Unix(),
	}

	token := jwt.NewWithClaims(
		jwt.SigningMethodHS256,
		claims,
	)

	return token.SignedString([]byte("secret-key"))
}
