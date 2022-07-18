package mqttcallbacks

import (
	"fmt"

	"github.com/cheetahfox/Iot-local-midware/health"
	mqtt "github.com/eclipse/paho.mqtt.golang"
)

var MessagePubHandler mqtt.MessageHandler = func(client mqtt.Client, msg mqtt.Message) {
	fmt.Printf("Message %s published on topic %s\n", msg.Payload(), msg.Topic())
}

// Set the Ready status for Kubernetes ready checks
var ConnectHandler mqtt.OnConnectHandler = func(client mqtt.Client) {
	fmt.Println("Connected to Broker")
	health.MqttReady = true
}

var ConnectionLostHandler mqtt.ConnectionLostHandler = func(client mqtt.Client, err error) {
	fmt.Printf("Connection Lost: %s\n", err.Error())
	health.MqttReady = false
}

var MessageSubHandler mqtt.MessageHandler = func(c mqtt.Client, msg mqtt.Message) {
	fmt.Println("Message %s received on topic %s\n", msg.Payload(), msg.Topic())
}
