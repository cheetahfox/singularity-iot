/*
Device specific functions for the Shelly25 Smart relay.

The Shelly 2.5 MQTT api is localted at https://shelly-api-docs.shelly.cloud/gen1/#shelly2-5-mqtt

Typical mqtt message format examples:

shellies/shellyswitch25-98CDAC38E9F5/input/0: 0
shellies/shellyswitch25-98CDAC38E9F5/input/1: 0
shellies/shellyswitch25-98CDAC38E9F5/temperature: 47.28
shellies/shellyswitch25-98CDAC38E9F5/temperature_f: 117.10
shellies/shellyswitch25-98CDAC38E9F5/overtemperature: 0
shellies/shellyswitch25-98CDAC38E9F5/temperature_status: Normal
shellies/shellyswitch25-98CDAC38E9F5/voltage: 122.26
shellies/shellyswitch25-98CDAC38E9F5/relay/0/power: 158.20
shellies/shellyswitch25-98CDAC38E9F5/relay/0/energy: 3572623

*/
package shelly

import (
	"errors"
	"fmt"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/cheetahfox/singularity-iot/database"
	mqtt "github.com/eclipse/paho.mqtt.golang"
	influxdb2 "github.com/influxdata/influxdb-client-go/v2"
)

type shelly25Device struct {
	macAddress string
	mode       string
	metric     map[string]float64
}

func makeShelly25dev() shelly25Device {
	var dev shelly25Device
	dev.metric = make(map[string]float64, 0)

	return dev
}

// Check to see if the metric string matches the valid options
func validateMetric(metric string) bool {
	verbs := []string{
		"power",
		"energy",
		"0",
		"1",
		"temperature",
		"temperature_f",
		"temperature_status",
		"overtemperature",
		"voltage",
	}

	for i := range verbs {
		if metric == verbs[i] {
			return true
		}
	}

	return false
}

/*
Here I parse through the Mtqq message and return a struct with the device details and the current metric/value
This is complicated by the fact that the length topic encodes what we data/case we have.
*/
func parseMessage25(msg mqtt.Message) (shelly25Device, error) {
	dev := makeShelly25dev()

	// All valid messages should be at least three values and less than 6 when split
	msgTopic := strings.Split(msg.Topic(), "/")
	if len(msgTopic) < 3 || len(msgTopic) >= 6 {
		fmt.Printf("Unable to split to Mqtt Topic: %s \n", msg.Topic())
		return dev, errors.New("unable to parse")
	}

	// This means we have a device metric
	if len(msgTopic) == 3 {
		if !validateMetric(msgTopic[3]) {
			errmessage := fmt.Sprintf("invalid metric : %s\n", msgTopic[3])
			return dev, errors.New(errmessage)
		}
		metric, err := strconv.ParseFloat(string(msg.Payload()), 64)
		if err != nil {
			return dev, err
		}
		dev.metric[msgTopic[3]] = metric
	}

	/*
		len 4 is if we have a relay or input message. Either relay/id on/of or input/id 0/1.
		I am going to simplify my life by making storing 0/1 for off/on. And seting the metric string to input-id/relay-id
	*/
	if len(msgTopic) == 4 {
		var metricName string
		switch msgTopic[3] {
		case "input":
			dev.mode = "input"
			metricName = fmt.Sprintf("input%s", msgTopic[4])
		case "relay":
			dev.mode = "relay"
			metricName = fmt.Sprintf("relay%s", msgTopic[4])
		case "roller":
			dev.mode = "roller"
			metricName = fmt.Sprintf("roller%s", msgTopic[4])
		}
		metric, err := strconv.ParseFloat(string(msg.Payload()), 64)
		if err != nil {
			return dev, err
		}
		if metricName != "" {
			dev.metric[metricName] = metric
		}

	}

	return dev, nil
}

func receiveMessage25(msg mqtt.Message) {
	fmt.Printf("Shelly 25 Device -  %s: %s\n", msg.Topic(), msg.Payload())

	v, _ := regexp.Compile(".+voltage$")
	topic := msg.Topic()

	dev, err := parseMessage25(msg)
	if err != nil {
		fmt.Println(err)
		return
	}

	if !validateMac(dev.macAddress) {
		return
	}

	if v.MatchString(topic) {
		payload := msg.Payload()
		currentTime := time.Now()

		voltage, err := strconv.ParseFloat(string(payload), 64)
		if err != nil {
			fmt.Println(err)
		}
		p := influxdb2.NewPointWithMeasurement("Shelly 2.5 Metrics")
		p.AddTag("Mac Address", dev.macAddress)
		p.AddField("Voltage", voltage)
		p.SetTime(currentTime)

		fmt.Println("Writing point: ", currentTime.Format(time.UnixDate))
		database.DbWrite.WritePoint(p)
	}
}
