package auth

import (
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v4"
)

func TokenNew(secret string) (string, error) {
	return jwt.
		NewWithClaims(jwt.SigningMethodHS256, jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(7 * 24 * time.Hour)),
		}).
		SignedString([]byte(secret))
}

func TokenParse(secret, tokenStr string) error {
	token, err := jwt.Parse(tokenStr, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(secret), nil
	})
	if err != nil {
		return fmt.Errorf("parsing token: %w", err)
	}

	if _, ok := token.Claims.(jwt.MapClaims); !(ok && token.Valid) {
		return fmt.Errorf("unauthorised")
	}

	return nil
}
