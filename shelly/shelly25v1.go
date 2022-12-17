package shelly

import (
	"encoding/json"

	"github.com/gofiber/fiber/v2"
)

/*
Returns a json formatted list of the Shelly 2.5 devices Mac Addresses
*/
func Api25ListDevsV1(c *fiber.Ctx) error {
	s, err := json.Marshal(shelly25Devs)
	if err != nil {
		return c.SendStatus(503)
	}

	c.Set("Content-Type", "application/json")
	return c.SendString(string(s))
}

/*
Sets the relay status:
topic: shellies/shellyswitch25-98CDAC38E9F5/relay/1/command msg: "on/off"
*/
func Api25RelayControl(c *fiber.Ctx) error {
	command := new(Shelly25Relay)

	err := c.BodyParser(command)
	if err != nil {
		return c.SendStatus(503)
	}

	return c.JSON(fiber.Map{"status": "success", "data": command})
}
