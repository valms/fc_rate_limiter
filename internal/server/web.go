package server

import (
	"context"
	"github.com/gofiber/fiber/v2"
	"github.com/valms/fc_rate_limiter/internal/cache"
	"github.com/valms/fc_rate_limiter/internal/config"
	"github.com/valms/fc_rate_limiter/internal/repository/redis"
	"github.com/valms/fc_rate_limiter/internal/server/middleware"
	"github.com/valms/fc_rate_limiter/internal/service"
	"net/http"
)

func SetupWebServer(loadConfig *config.Config) *fiber.App {
	app := fiber.New()
	redisAddr := loadConfig.Redis.Host + ":" + loadConfig.Redis.Port

	redisClient, _ := cache.NewRedisClient(redisAddr, loadConfig.Redis.Port)

	operations, _ := redis.NewRedisOperations(context.Background(), redisClient)

	rate, _ := service.NewLimiterService(operations, loadConfig)

	app.Use(middleware.Limiter(rate))

	app.Get("/", func(c *fiber.Ctx) error {
		return c.SendStatus(http.StatusOK)
	})

	return app
}
