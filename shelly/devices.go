/*
Shelly Specific functions.
*/

package shelly

import (
	"fmt"
	"regexp"
	"time"

	mqtt "github.com/eclipse/paho.mqtt.golang"
)

// shellies/shellyswitch25-98CDAC38E9F5/temperature: 45.90
var shellyTempHandler mqtt.MessageHandler = func(client mqtt.Client, msg mqtt.Message) {
	fmt.Printf("%s : %s\n", msg.Topic(), msg.Payload())
}

// shellies/shellyswitch25-98CDAC38E9F5/relay/0/power: 117.89
var shellyPowerHandler mqtt.MessageHandler = func(client mqtt.Client, msg mqtt.Message) {
	fmt.Printf("%s : %s\n", msg.Topic(), msg.Payload())
}

func ReceiveMessage(msg mqtt.Message) {
	// Device specific callbacks
	shelly25Re, _ := regexp.Compile("shellies/shellyswitch25-.+$")
	shelly15Re, _ := regexp.Compile("shellies/shellyswitch15-.+$")
	shellyAnnouce, _ := regexp.Compile("shellies/announce.+$")

	switch true {
	case shelly15Re.MatchString(msg.Topic()):
		receiveMessage15(msg)
	case shelly25Re.MatchString(msg.Topic()):
		receiveMessage25(msg)
	case shellyAnnouce.MatchString(msg.Topic()):
		receiveAnnounce(msg)
	default:
		now := time.Now()
		fmt.Println("unknown Shelly message: not processed")
		fmt.Printf("Published on Topic: %s value: %s @ %s \n", msg.Topic(), msg.Payload(), now.Format(time.UnixDate))
	}
}

// currently don't validate the shelly mac address it's not a standard format - IMPLEMENT THIS
func validateMac(mac string) bool {
	return true
}

/*
Published on Topic: shellies/announce  value: {"id":"shellyswitch25-98CDAC38E9F5","model":"SHSW-25","mac":"98CDAC38E9F5","ip":"192.168.76.119","new_fw":true,"fw_ver":"20220209-093016/v1.11.8-g8c7bb8d","mode":"relay"}
Shelly 25 Device -  shellies/shellyswitch25-98CDAC38E9F5/announce: {"id":"shellyswitch25-98CDAC38E9F5","model":"SHSW-25","mac":"98CDAC38E9F5","ip":"192.168.76.119","new_fw":true,"fw_ver":"20220209-093016/v1.11.8-g8c7bb8d","mode":"relay"}
This function will register the new shelly device for now it just logs the message to stdout
*/
func receiveAnnounce(msg mqtt.Message) {
	fmt.Println("Shelly device annoucement")
	fmt.Println(msg.Payload())
}

/*
Generic shelly device temp func

Example mqtt output.

Shelly 25 Device -  shellies/shellyswitch25-98CDAC38E9F5/temperature: 45.90

*/
func shellyTempSub(client mqtt.Client, macAddr string) {
	topic := "shellies/shellyswitch25-" + macAddr + "/temperature"
	client.Subscribe(topic, 0, shellyTempHandler)
}

// shellies/shellyswitch25-98CDAC38E9F5/relay/0/power: 117.89
func shellyPowerSub(client mqtt.Client, macAddr string, relay string) {
	topic := "hellies/shellyswitch25" + macAddr + "/realy/" + relay + "/power"
	client.Subscribe(topic, 0, shellyPowerHandler)
}
