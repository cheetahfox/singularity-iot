package shelly

import (
	"fmt"
	"regexp"

	mqtt "github.com/eclipse/paho.mqtt.golang"
)

func ReceiveMessage(msg mqtt.Message) {
	// Shellyswitch25
	shelly25Re, _ := regexp.Compile("shellies/shellyswitch25-.+$")
	shelly15Re, _ := regexp.Compile("shellies/shellyswitch15-.+$")

	switch true {
	case shelly15Re.MatchString(msg.Topic()):
		receiveMessage15(msg)
	case shelly25Re.MatchString(msg.Topic()):
		receiveMessage25(msg)
	default:
		fmt.Println("unknown Shelly message: not processed")
		fmt.Printf("Published on Topic: %s  value: %s\n", msg.Topic(), msg.Payload())
	}
}
