package server

import (
	"Rate_Limiter/internal/cache"
	"Rate_Limiter/internal/config"
	"Rate_Limiter/internal/repository/redis"
	"Rate_Limiter/internal/server/middleware"
	"Rate_Limiter/internal/service"
	"context"
	"github.com/gofiber/fiber/v2"
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
