/*

Device specific functions for the Shelly25 Smart relay.

The Shelly 2.5 MQTT api is localted at https://shelly-api-docs.shelly.cloud/gen1/#shelly2-5-mqtt

In


*/
package shelly

import (
	"fmt"
	"regexp"
	"strconv"
	"time"

	"github.com/cheetahfox/Iot-local-midware/database"
	mqtt "github.com/eclipse/paho.mqtt.golang"
	influxdb2 "github.com/influxdata/influxdb-client-go/v2"
)

func receiveMessage25(msg mqtt.Message) {
	fmt.Printf("Shelly 25 Device -  %s: %s\n", msg.Topic(), msg.Payload())

	v, _ := regexp.Compile(".+voltage$")
	topic := msg.Topic()

	if v.MatchString(topic) {
		payload := msg.Payload()

		voltage, err := strconv.ParseFloat(string(payload), 64)
		if err != nil {
			fmt.Println(err)
		}
		p := influxdb2.NewPointWithMeasurement("Shelly 2.5 Metrics")
		p.AddTag("Mac Address", "98CDAC38E9F5")
		p.AddField("Voltage", voltage)
		p.SetTime(time.Now())

		fmt.Println("Writing point")
		database.DbWrite.WritePoint(p)
	}
}
