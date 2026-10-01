package middleware

import (
	"log"
	"time"

	"github.com/gofiber/fiber/v2"
)

func Logger() fiber.Handler {

	return func(c *fiber.Ctx) error {

		start := time.Now()

		err := c.Next()

		duration := time.Since(start)

		log.Printf(
			"[%s] %d %s %s (%v)",
			time.Now().Format("15:04:05"),
			c.Response().StatusCode(),
			c.Method(),
			c.Path(),
			duration,
		)

		return err
	}
}