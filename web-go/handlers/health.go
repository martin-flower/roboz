package handlers

import (
	"time"

	"github.com/gofiber/fiber/v2"
	"go.uber.org/zap"
)

// Health
// @Summary service health status
// @Description returns 200 if service is healthy (shallow health)
// @Tags health check
// @Produce json
// @Success 200
// @Failure 500
// @Router /health [get]
func Health(c *fiber.Ctx) error {
	zap.S().Infof("roboz health ok - %d", time.Now().UnixMilli())
	return c.Status(fiber.StatusOK).SendString("ok health " + time.Now().Format(time.RFC3339))
}
