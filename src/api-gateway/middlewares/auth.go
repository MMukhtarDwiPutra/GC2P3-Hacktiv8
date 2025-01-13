package middlewares

import (
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo/v4"
)

func RequireAuth(next echo.HandlerFunc) echo.HandlerFunc {
	SecretKey := os.Getenv("LOGIN_SECRET_KEY")

	return func(c echo.Context) error {
		// Retrieve token from the Authorization header.
		authHeader := c.Request().Header.Get("Authorization")
		if authHeader == "" {
			return c.JSON(http.StatusUnauthorized, map[string]string{"message": "Missing Authorization header"})
		}

		// Extract the token from the "Bearer <token>" format.
		tokenString := strings.TrimPrefix(authHeader, "Bearer ")
		if tokenString == authHeader {
			return c.JSON(http.StatusUnauthorized, map[string]string{"message": "Invalid token format"})
		}

		// Parse the token using the secret key.
		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			// Validate the signing method.
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, echo.NewHTTPError(http.StatusUnauthorized, "Invalid signing method")
			}
			return []byte(SecretKey), nil
		})
		if err != nil {
			return c.JSON(http.StatusUnauthorized, map[string]string{"message": "Invalid token"})
		}

		// Validate the claims within the token.
		if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
			// Check the token expiration time.
			exp, ok := claims["exp"].(float64)
			if !ok {
				return c.JSON(http.StatusUnauthorized, map[string]string{"message": "Invalid expiration time in token"})
			}

			// Ensure the token is not expired.
			if time.Unix(int64(exp), 0).Before(time.Now()) {
				return c.JSON(http.StatusUnauthorized, map[string]string{"message": "Token has expired"})
			}

			// Extract and set the user_id claim to the context.
			userID, ok := claims["user_id"].(string)
			if !ok {
				return c.JSON(http.StatusUnauthorized, map[string]string{"message": "Missing or invalid user_id in token claims"})
			}
			c.Set("user", claims)    // Set the full claims in the context for future use.
			c.Set("user_id", userID) // Also set the user_id separately for convenience.

			// Proceed to the next handler.
			return next(c)
		}

		// If claims are invalid.
		return c.JSON(http.StatusUnauthorized, map[string]string{"message": "Invalid token claims"})
	}
}
