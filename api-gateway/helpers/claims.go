package helpers

import (
	"api-gateway/dto"
	"fmt"

	"github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo/v4"
)

func GetClaims(c echo.Context) (dto.Claims, error) {
	claimsTmp := c.Get("user")
	if claimsTmp == nil {
		return dto.Claims{}, fmt.Errorf("Failed to fetch user claims from JWT")
	}

	claims := claimsTmp.(jwt.MapClaims)
	return dto.Claims{
		UserID: uint(claims["user_id"].(float64)),
	}, nil
}
