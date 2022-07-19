package mqttcallbacks

import (
	"fmt"
	"regexp"

	"github.com/cheetahfox/Iot-local-midware/config"
	"github.com/cheetahfox/Iot-local-midware/health"
	shelly "github.com/cheetahfox/Iot-local-midware/shelly"
	mqtt "github.com/eclipse/paho.mqtt.golang"
)

var MessagePubHandler mqtt.MessageHandler = func(client mqtt.Client, msg mqtt.Message) {
	fmt.Printf("Message %s published on topic %s\n", msg.Payload(), msg.Topic())
}

// Set the Ready status for Kubernetes ready checks
var ConnectHandler mqtt.OnConnectHandler = func(client mqtt.Client) {
	ops := client.OptionsReader()
	servers := ops.Servers()
	for index := range servers {
		fmt.Printf("Connected to Broker %s\n", servers[index].Hostname())
	}

	health.MqttReady = true
}

var ConnectionLostHandler mqtt.ConnectionLostHandler = func(client mqtt.Client, err error) {
	fmt.Printf("Connection Lost: %s\n", err.Error())
	health.MqttReady = false
}

// This will process all incoming subscribed messages. Here we will call device specific functions
var MessageSubHandler mqtt.MessageHandler = func(c mqtt.Client, msg mqtt.Message) {
	// Shelly Devices
	shellyRe, _ := regexp.Compile("shellies/.+$")

	switch true {
	case shellyRe.MatchString(msg.Topic()):
		shelly.ReceiveMessage(msg)
	default:
		fmt.Printf("Message %s received on topic %s\n", msg.Payload(), msg.Topic())
	}
}

// Need to set the defaults here in this package to keep from having a import cycle problems.
func SetDefaultCallbacks(c *config.Configuration) {
	c.Options.SetDefaultPublishHandler(MessagePubHandler)
	c.Options.OnConnect = ConnectHandler
	c.Options.OnConnectionLost = ConnectionLostHandler
}
