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
