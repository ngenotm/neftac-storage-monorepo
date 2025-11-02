package api

import (
	"strings"
	"neftac/storage/internal/auth"
	"neftac/storage/internal/policy"
	"github.com/dgrijalva/jwt-go"
	"github.com/gofiber/fiber/v2"
)

func Auth(c *fiber.Ctx) error {
	h := c.Get("Authorization")
	if !strings.HasPrefix(h, "Bearer ") { return c.Status(401).JSON(fiber.Map{"error": "unauthorized"}) }
	t := strings.TrimPrefix(h, "Bearer ")
	claims := &auth.Claims{}
	_, err := jwt.ParseWithClaims(t, claims, func(*jwt.Token) (interface{}, error) { return auth.Secret, nil })
	if err != nil { return c.Status(401).JSON(fiber.Map{"error": "invalid token"}) }
	c.Locals("user", claims.UserID)
	c.Locals("role", claims.Role)
	return c.Next()
}

func Policy(act string) fiber.Handler {
	return func(c *fiber.Ctx) error {
		user := c.Locals("user").(string)
		b := c.Params("bucket")
		p := c.Params("path")
		if !policy.Allow(user, b+":"+p, act) { return c.Status(403).JSON(fiber.Map{"error": "forbidden"}) }
		return c.Next()
	}
}
