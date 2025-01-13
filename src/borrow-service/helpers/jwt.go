package helpers

import (
	"borrow-service/dto"
	"fmt"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo/v4"
)

func SignJwtForGrpc() (string, error) {
	secret := os.Getenv("SERVICES_JWT_SECRET")

	claims := &jwt.RegisteredClaims{
		ExpiresAt: jwt.NewNumericDate(time.Now().AddDate(0, 0, 1)),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(secret))
	if err != nil {
		return "", fmt.Errorf("failed to generate JWT: %v", err)
	}
	return tokenString, nil
}

func GetClaims(c echo.Context) (dto.Claims, error) {
	claimsTmp := c.Get("user")
	if claimsTmp == nil {
		return dto.Claims{}, fmt.Errorf("Failed to fetch user claims from JWT")
	}

	// Assert claims to jwt.MapClaims
	claims, ok := claimsTmp.(jwt.MapClaims)
	if !ok {
		return dto.Claims{}, fmt.Errorf("Failed to assert user claims")
	}

	// Assert user_id to string
	userID, ok := claims["user_id"].(string)
	if !ok {
		return dto.Claims{}, fmt.Errorf("Failed to fetch user_id from claims")
	}

	return dto.Claims{
		UserID: userID,
	}, nil
}
