package middlewares

import (
	"github.com/fakhrads/golang-test/utils"
	"github.com/gofiber/fiber/v2"
)

func KeyValidator() fiber.Handler {
	return func(c *fiber.Ctx) error {
		apiKey := c.Get("key")
		expectedKey := utils.GetEnv("API_KEY")

		if apiKey == "" {
			return c.Status(403).JSON(fiber.Map{"error": "API key is missing."})
		}

		if apiKey != expectedKey {
			return c.Status(401).JSON(fiber.Map{"error": "Invalid API key."})
		}

		return c.Next()
	}
}
