package health

import (
	"fmt"
	"time"

	"github.com/gofiber/fiber/v2"
)

var MqttReady, InfluxReady bool
var LastRecieved time.Time
var PointsWritten int64

func init() {
	MqttReady = false
	InfluxReady = false
	LastRecieved = time.Now()
	go receiveCheck()
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

/*
Check to see if we are getting Mqtt messages if we don't after 5 minutes we set not ready
I am doing this since I have seen the OnConnectHandler doesn't always reconnect
*/
func receiveCheck() {
	ticker := time.NewTicker(time.Second * time.Duration(15))
	for range ticker.C {
		now := time.Now()
		if now.Sub(LastRecieved) >= (time.Second * time.Duration(300)) {
			MqttReady = false
			fmt.Printf("300 seconds or more since mqtt message recieved: marking not ready: %s\n", now.Format(time.UnixDate))
		}
	}
}
