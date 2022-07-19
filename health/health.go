package health

import (
	"github.com/gofiber/fiber/v2"
)

var MqttReady, InfluxReady bool

func init() {
	MqttReady = false
	InfluxReady = false
}

func GetHealthz(c *fiber.Ctx) error {
	// return &fiber.Error{}
	return c.SendStatus(200)
}

func GetReadyz(c *fiber.Ctx) error {
	if !MqttReady || !InfluxReady {
		return c.SendStatus(503)
	}
	return c.SendStatus(200)
}
