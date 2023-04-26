package auth

import "github.com/gofiber/fiber/v2"

// CreateHttpOnlyCookie is a function that creates a cookie
// that is only accessible through HTTP requests
func CreateHttpOnlyCookie(name string, value string, c *fiber.Ctx) {
	c.Cookie(&fiber.Cookie{
		Name:     name,
		Value:    value,
		HTTPOnly: true,
	})
}