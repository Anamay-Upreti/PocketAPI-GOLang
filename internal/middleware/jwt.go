package middleware

import (
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"

	"pocketapi/internal/config"
	"pocketapi/internal/utils"
)

func Protected() fiber.Handler {

	return func(c *fiber.Ctx) error {

		authHeader := c.Get("Authorization")

		if authHeader == "" {
			return utils.Error(c, 401, "Authorization header missing")
		}

		if !strings.HasPrefix(authHeader, "Bearer ") {
			return utils.Error(c, 401, "Invalid authorization format")
		}

		tokenString := strings.TrimPrefix(authHeader, "Bearer ")

		token, err := jwt.ParseWithClaims(
			tokenString,
			&utils.JWTClaims{},
			func(token *jwt.Token) (interface{}, error) {
				return []byte(config.AppConfig.JWTSecret), nil
			},
		)

		if err != nil || !token.Valid {
			return utils.Error(c, 401, "Invalid or expired token")
		}

		claims := token.Claims.(*utils.JWTClaims)

		c.Locals("userID", claims.UserID)
		c.Locals("email", claims.Email)
		c.Locals("role", claims.Role)

		return c.Next()
	}
}