package middleware

import (
	"fmt"
	"strings"

	"github.com/HouseCham/VetMate/config"
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt"
)

var Config *config.Config

func ShareConfigFile(config *config.Config) {
	Config = config
}

func JwtMiddleware() fiber.Handler {
	return func(c *fiber.Ctx) error {
		// Get JWT token from Authorization header
		authHeader := c.Get("Authorization")
		tokenString := strings.Replace(authHeader, "Bearer ", "", 1)

		// Parse JWT token
		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			// Check if signing method is HMAC
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("invalid signing method")
			}
			// Return secret key for HMAC validation
			return []byte(Config.DevConfiguration.Jwt.Secret), nil
		})
		if err != nil {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"message": "Unauthorized",
			})
		}

		// Check if token is valid
		if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
			// Set authenticated user ID in request context
			c.Locals("userId", claims["sub"])
			return c.Next()
		}

		// Token is invalid
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"message": "Unauthorized",
		})
	}
}
