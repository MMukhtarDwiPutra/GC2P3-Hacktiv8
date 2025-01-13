package helpers

import (
	"api-gateway/dto"
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo/v4"
)

func SignNewJWT(c echo.Context, user dto.User) error {
	claims := jwt.MapClaims{
		"exp":      time.Now().Add(4 * time.Hour).Unix(),
		"user_id":  user.UserID,
		"username": user.Username,
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(os.Getenv("JWT_SECRET")))
	if err != nil {
		return fmt.Errorf("Error sign token")
	}

	cookie := new(http.Cookie)
	cookie.Name = "Authorization"
	cookie.HttpOnly = true
	cookie.Path = "/"
	cookie.Value = tokenString
	cookie.SameSite = http.SameSiteLaxMode
	cookie.Expires = time.Now().Add(2 * time.Hour)
	c.SetCookie(cookie)

	return nil
}

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
