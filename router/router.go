package router

import (
	"github.com/cheetahfox/singularity-iot/health"
	"github.com/cheetahfox/singularity-iot/shelly"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
)

func SetupRoutes(app *fiber.App) {
	app.Get("/", func(c *fiber.Ctx) error {
		return c.SendString("Singularity Running") // send text
	})

	app.Get("/healthz", health.GetHealthz)
	app.Get("/readyz", health.GetReadyz)

	// Shelly 2.5 v1 Web API
	shelly25 := app.Group("/api/shelly2.5/v1", logger.New())
	shelly25.Get("list", shelly.Api25ListDevsV1)
	shelly25.Put("relay", shelly.Api25RelayControl)

	//api := app.Group("/api/v1/", logger.New())
	//api.Get("growlights", growlightv1.GetGrowLights)
}
