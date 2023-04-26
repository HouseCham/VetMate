package auth

import (
	"errors"

	"github.com/gofiber/fiber/v2"
)

// CreateHttpOnlyCookie is a function that creates a cookie
// that is only accessible through HTTP requests
func CreateJwtCookie(value string, c *fiber.Ctx) {
	c.Cookie(&fiber.Cookie{
		Name:     Config.DevConfiguration.Jwt.CookieName,
		Value:    value,
		HTTPOnly: true,
		Path:     "/",
	})
}

// GetHttpOnlyCookie is a function that gets the value of a cookie
// that is only accessible through HTTP requests
func GetJwtCookie(c *fiber.Ctx) (string, error) {
	cookie := c.Cookies(Config.DevConfiguration.Jwt.CookieName)
	if cookie == "" {
		return "", errors.New("cookie not found")
	}
	return cookie, nil
}

// DeleteHttpOnlyCookie is a function that deletes a cookie
// that is only accessible through HTTP requests
func DeleteJwtCookie(c *fiber.Ctx) {
	c.Cookie(&fiber.Cookie{
		Name:     Config.DevConfiguration.Jwt.CookieName,
		Value:    "",
		HTTPOnly: true,
		MaxAge:   -1,
		Path:     "/",
	})
}
