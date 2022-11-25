package router

import (
	"github.com/cheetahfox/singularity-iot/health"
	"github.com/gofiber/fiber/v2"
)

func SetupRoutes(app *fiber.App) {
	app.Get("/", func(c *fiber.Ctx) error {
		return c.SendString("Iot Middleware running") // send text
	})

	app.Get("/healthz", health.GetHealthz)
	app.Get("/readyz", health.GetReadyz)

	//api := app.Group("/api/v1/", logger.New())
	//api.Get("growlights", growlightv1.GetGrowLights)
}
