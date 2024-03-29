/*
Shelly generic functions for subscribing to topics or directing to device specific.
*/

package shelly

import (
	"fmt"
	"regexp"
	"time"

	"github.com/cheetahfox/singularity-iot/health"
	mqtt "github.com/eclipse/paho.mqtt.golang"
)

/*
Since this is a generic temp call back we have to route the messages to the specific shelly device code.
example msg : shellies/shellyswitch25-98CDAC38E9F5/temperature: 45.90
*/
var shellyTempHandler mqtt.MessageHandler = func(client mqtt.Client, msg mqtt.Message) {
	health.LastRecieved = time.Now()

	// Device specific callbacks
	shelly25Re, _ := regexp.Compile("shellies/shellyswitch25-.+$")
	shelly15Re, _ := regexp.Compile("shellies/shellyswitch15-.+$")

	switch true {
	case shelly15Re.MatchString(msg.Topic()):
		rcv15Temp(msg)
	case shelly25Re.MatchString(msg.Topic()):
		err := rcv25Temp(msg)
		if err != nil {
			fmt.Println(err)
		}
	default:
		fmt.Printf("shellyTempHandler ---> Unknown %s : %s\n", msg.Topic(), msg.Payload())
	}
}

// shellies/shellyswitch25-98CDAC38E9F5/relay/0/power: 117.89
var shellyPowerHandler mqtt.MessageHandler = func(client mqtt.Client, msg mqtt.Message) {
	health.LastRecieved = time.Now()

	shelly25Re, _ := regexp.Compile("shellies/shellyswitch25-.+$")
	shelly15Re, _ := regexp.Compile("shellies/shellyswitch15-.+$")

	switch true {
	case shelly15Re.MatchString(msg.Topic()):
		err := rcv15Power(msg)
		if err != nil {
			fmt.Println(err)
		}
	case shelly25Re.MatchString(msg.Topic()):
		err := rcv25Power(msg)
		if err != nil {
			fmt.Println(err)
		}
	default:
		fmt.Printf("shellyTempHandler ---> Unknown %s : %s\n", msg.Topic(), msg.Payload())
	}
}

// shellies/shellyswitch25-98CDAC38E9F5/voltage: 122.26
var shellyVoltageHandler mqtt.MessageHandler = func(client mqtt.Client, msg mqtt.Message) {
	health.LastRecieved = time.Now()

	shelly25Re, _ := regexp.Compile("shellies/shellyswitch25-.+$")
	shelly15Re, _ := regexp.Compile("shellies/shellyswitch15-.+$")

	switch true {
	case shelly15Re.MatchString(msg.Topic()):
		err := rcv15Voltage(msg)
		if err != nil {
			fmt.Println(err)
		}
	case shelly25Re.MatchString(msg.Topic()):
		err := rcv25Voltage(msg)
		if err != nil {
			fmt.Println(err)
		}
	default:
		fmt.Printf("shellyVoltageHandler ---> Unknown %s : %s\n", msg.Topic(), msg.Payload())
	}
}

// shellies/shellyswitch25-98CDAC38E9F5/relay/0/energy: 3572623
var shellyEnergyHandler mqtt.MessageHandler = func(c mqtt.Client, m mqtt.Message) {
	health.LastRecieved = time.Now()

	shelly25Re, _ := regexp.Compile("shellies/shellyswitch25-.+$")
	shelly15Re, _ := regexp.Compile("shellies/shellyswitch15-.+$")

	switch true {
	case shelly15Re.MatchString(m.Topic()):
		err := rcv15Energy(m)
		if err != nil {
			fmt.Println(err)
		}
	case shelly25Re.MatchString(m.Topic()):
		err := rcv25Energy(m)
		if err != nil {
			fmt.Println(err)
		}
	default:
		fmt.Printf("shellyTempHandler ---> Unknown %s : %s\n", m.Topic(), m.Payload())
	}
}

/*
Published on Topic: shellies/announce  value: {"id":"shellyswitch25-98CDAC38E9F5","model":"SHSW-25","mac":"98CDAC38E9F5","ip":"192.168.76.119","new_fw":true,"fw_ver":"20220209-093016/v1.11.8-g8c7bb8d","mode":"relay"}
Shelly 25 Device -  shellies/shellyswitch25-98CDAC38E9F5/announce: {"id":"shellyswitch25-98CDAC38E9F5","model":"SHSW-25","mac":"98CDAC38E9F5","ip":"192.168.76.119","new_fw":true,"fw_ver":"20220209-093016/v1.11.8-g8c7bb8d","mode":"relay"}
This function will register the new shelly device for now it just logs the message to stdout
*/

/*
Generic shelly device temp func

Example mqtt output.

Shelly 25 Device -  shellies/shellyswitch25-98CDAC38E9F5/temperature: 45.90

*/
func shelly25TempSub(client mqtt.Client, macAddr string) {
	topic := "shellies/shellyswitch25-" + macAddr + "/temperature"
	client.Subscribe(topic, 0, shellyTempHandler)
}

// shellies/shellyswitch25-98CDAC38E9F5/relay/0/power: 117.89
func shelly25PowerSub(client mqtt.Client, macAddr string, relay string) {
	topic := "shellies/shellyswitch25-" + macAddr + "/relay/" + relay + "/power"
	client.Subscribe(topic, 0, shellyPowerHandler)
}

func shelly25EnergySub(c mqtt.Client, macAddr string, relay string) {
	topic := "shellies/shellyswitch25-" + macAddr + "/relay/" + relay + "/energy"
	c.Subscribe(topic, 0, shellyEnergyHandler)
}

// shellies/shellyswitch25-98CDAC38E9F5/voltage: 123.29
func shelly25VotlageSub(client mqtt.Client, macAddr string) {
	topic := "shellies/shellyswitch25-" + macAddr + "/voltage"
	client.Subscribe(topic, 0, shellyVoltageHandler)
}
