package middleware

import (
	"github.com/gofiber/fiber/v2"
	"github.com/valms/fc_rate_limiter/internal/ratelimiter"
)

type Config struct {
}

func Limiter(rate *ratelimiter.RateLimiter) fiber.Handler {
	return func(c *fiber.Ctx) error {
		ip := c.IP()
		apikey := c.Get("API_KEY")

		var blocked bool

		if apikey != "" {
			blocked, _ = rate.IsRateLimitByToken(apikey)
		} else {
			blocked, _ = rate.IsRateLimitByIP(ip)
		}

		if blocked {
			return c.Status(429).JSON(&fiber.Map{
				"message": "you have reached the maximum number of requests or actions allowed within a certain time frame",
			})
		}

		return c.Next()

	}
}
