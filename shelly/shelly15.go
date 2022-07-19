package shelly

import (
	"fmt"

	mqtt "github.com/eclipse/paho.mqtt.golang"
)

func receiveMessage15(msg mqtt.Message) {
	fmt.Printf("Shelly 25 Device -  %s: %s\n", msg.Topic(), msg.Payload())
}
