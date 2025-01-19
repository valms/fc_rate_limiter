package middleware

import (
	"Rate_Limiter/internal/service"
	"github.com/gofiber/fiber/v2"
)

type Config struct {
}

func Limiter(rate *service.RateLimiter) fiber.Handler {
	return func(c *fiber.Ctx) error {
		ip := c.IP()

		ipBlocked, _ := rate.IsRateLimitByIP(ip)

		if ipBlocked {
			return c.Status(429).JSON(&fiber.Map{
				"message": "you have reached the maximum number of requests or actions allowed within a certain time frame",
			})
		}

		return c.Next()

	}
}
