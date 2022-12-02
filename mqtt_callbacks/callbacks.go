package mqttcallbacks

import (
	"fmt"
	"regexp"
	"time"

	"github.com/cheetahfox/singularity-iot/config"
	"github.com/cheetahfox/singularity-iot/health"
	shelly "github.com/cheetahfox/singularity-iot/shelly"
	mqtt "github.com/eclipse/paho.mqtt.golang"
)

var lastRecieved time.Time

func init() {
	lastRecieved = time.Now()
	go receiveCheck()
}

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

// This will process all incoming subscribed messages. Here we will call device family specific functions
var MessageSubHandler mqtt.MessageHandler = func(c mqtt.Client, msg mqtt.Message) {
	lastRecieved = time.Now()
	// Shelly Devices
	shellyRe, _ := regexp.Compile("shellies/.+$")

	switch true {
	case shellyRe.MatchString(msg.Topic()):
		shelly.ReceiveMessage(msg)
	default:
		//fmt.Printf("Message %s received on topic %s\n", msg.Payload(), msg.Topic())
	}
}

// Need to set the defaults here in this package to keep from having a import cycle problems.
func SetDefaultCallbacks(c *config.Configuration) {
	c.Options.SetDefaultPublishHandler(MessagePubHandler)
	c.Options.OnConnect = ConnectHandler
	c.Options.OnConnectionLost = ConnectionLostHandler
}

/*
Check to see if we are getting Mqtt messages if we don't after 5 minutes we set not ready
I am doing this since I have seen the OnConnectHandler doesn't always reconnect
*/
func receiveCheck() {
	ticker := time.NewTicker(time.Second * time.Duration(15))
	for range ticker.C {
		now := time.Now()
		if now.Sub(lastRecieved) >= (time.Second * time.Duration(300)) {
			health.MqttReady = false
			fmt.Printf("300 seconds or more since mqtt message recieved: marking not ready: %s", now.Format(time.UnixDate))
		}
	}
}
