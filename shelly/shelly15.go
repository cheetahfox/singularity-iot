package shelly

import (
	"fmt"

	mqtt "github.com/eclipse/paho.mqtt.golang"
)

// rcv Temp stub func
func rcv15Temp(msg mqtt.Message) error {
	fmt.Printf("Shelly15 Device -  %s: %s\n", msg.Topic(), msg.Payload())

	return nil
}

// rcv Power stub func
func rcv15Power(msg mqtt.Message) error {
	fmt.Printf("Shelly15 Device -  %s: %s\n", msg.Topic(), msg.Payload())

	return nil
}

// rcv Energy stub func
func rcv15Energy(msg mqtt.Message) error {
	fmt.Printf("Shelly15 Device -  %s: %s\n", msg.Topic(), msg.Payload())

	return nil
}

// rcv Voltage stub func
func rcv15Voltage(msg mqtt.Message) error {
	fmt.Printf("Shelly15 Device -  %s: %s\n", msg.Topic(), msg.Payload())

	return nil
}
