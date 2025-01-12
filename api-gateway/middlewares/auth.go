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
		// Mendapatkan token dari header Authorization.
		authHeader := c.Request().Header.Get("Authorization")
		if authHeader == "" {
			// Mengembalikan respons Unauthorized jika header Authorization tidak ditemukan.
			return c.JSON(http.StatusUnauthorized, map[string]string{"message": "Missing Authorization header"})
		}

		// Memisahkan token dari format "Bearer <token>".
		tokenString := strings.TrimPrefix(authHeader, "Bearer ")
		if tokenString == authHeader {
			// Mengembalikan respons Unauthorized jika format token tidak valid.
			return c.JSON(http.StatusUnauthorized, map[string]string{"message": "Invalid token format"})
		}

		// Memparsing token menggunakan kunci rahasia.
		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			// Memvalidasi metode signing token.
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, echo.NewHTTPError(http.StatusUnauthorized, "Invalid signing method")
			}
			return []byte(SecretKey), nil
		})
		if err != nil {
			// Mengembalikan respons Unauthorized jika token tidak valid.
			return c.JSON(http.StatusUnauthorized, map[string]string{"message": "Invalid token"})
		}

		// Memvalidasi klaim (claims) dalam token.
		if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
			// Memeriksa waktu kedaluwarsa (expiration time) token.
			exp, ok := claims["exp"].(float64)
			if !ok {
				return c.JSON(http.StatusUnauthorized, map[string]string{"message": "Invalid expiration time"})
			}

			// Memastikan token belum kedaluwarsa.
			if time.Unix(int64(exp), 0).Before(time.Now()) {
				return c.JSON(http.StatusUnauthorized, map[string]string{"message": "Token has expired"})
			}

			if userID, ok := claims["user_id"].(float64); ok {
				c.Set("user_id", int(userID)) // Mengonversi ke int sebelum menyimpannya.
			}

			// Melanjutkan ke handler berikutnya jika token valid.
			return next(c)
		}

		// Mengembalikan respons Unauthorized jika klaim dalam token tidak valid.
		return c.JSON(http.StatusUnauthorized, map[string]string{"message": "Invalid token claims"})
	}
}
