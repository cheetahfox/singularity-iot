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
	"strconv"
	"strings"
	"time"

	"github.com/cheetahfox/singularity-iot/config"
	"github.com/cheetahfox/singularity-iot/database"
	"github.com/cheetahfox/singularity-iot/health"
	mqtt "github.com/eclipse/paho.mqtt.golang"
	influxdb2 "github.com/influxdata/influxdb-client-go/v2"
)

type shelly25Data struct {
	macAddress string
	relay      string
	metric     map[string]float64
}

func makeShelly25data() shelly25Data {
	var data shelly25Data
	data.metric = make(map[string]float64, 0)

	return data
}

// Init a new Shelly25
func InitShelly25dev(client mqtt.Client, device config.Iotdevices) error {
	shelly25TempSub(client, device.Maddr)
	shelly25PowerSub(client, device.Maddr, "0")
	shelly25PowerSub(client, device.Maddr, "1")
	shelly25VotlageSub(client, device.Maddr)
	shelly25EnergySub(client, device.Maddr, "0")
	shelly25EnergySub(client, device.Maddr, "1")
	fmt.Println("Shelly 25 device: " + device.Maddr + " Init Complete")

	return nil
}

/*
Receive Shelly 2.5 Temp
*/
func rcv25Temp(msg mqtt.Message) error {
	dp := makeShelly25data()

	// Get the Mac Address by removing the front and back of the expected topic string
	preformatMac := strings.TrimPrefix(msg.Topic(), "shellies/shellyswitch25-")
	dp.macAddress = strings.TrimSuffix(preformatMac, "/temperature")
	metric := "Temperature"

	data, err := strconv.ParseFloat(string(msg.Payload()), 64)
	if err != nil {
		return err
	}

	dp.metric[metric] = data
	write25point(dp, metric)

	return nil
}

/*
Receive Shelly 2.5 Voltage
*/
func rcv25Voltage(m mqtt.Message) error {
	dp := makeShelly25data()
	metric := "Voltage"

	dp, err := parseVals(dp, metric, m)
	if err != nil {
		return err
	}
	write25point(dp, metric)

	return nil
}

/*
Receive Shelly 2.5 current power in watts
Topic/Metric output shellies/shellyswitch25-98CDAC38E9F5/relay/0/power: 158.20
*/
func rcv25Power(msg mqtt.Message) error {
	dp := makeShelly25data()
	metric := "Power"

	dp, err := relayParseVals(dp, metric, msg)
	if err != nil {
		return err
	}
	write25point(dp, metric)

	return nil
}

/*
Recieve Shelly 2.5 energy
This is in a total of watt minutes since the device was powered on/rebooted.
This is kinda of an odd metric to work with but it's rather accurate and the best way to figure out
long term power usage.
*/
func rcv25Energy(msg mqtt.Message) error {
	dp := makeShelly25data()
	metric := "Energy"

	dp, err := relayParseVals(dp, metric, msg)
	if err != nil {
		return err
	}
	write25point(dp, metric)

	return nil
}

// Parse topic and fill in the metric value
func parseVals(dp shelly25Data, metricName string, m mqtt.Message) (shelly25Data, error) {
	// Basic format checking and filling in values if we have correctly formated Topics
	t := strings.Split(m.Topic(), "/")

	if len(t) != 3 {
		err := errors.New("malformed topic: " + m.Topic())
		return dp, err
	}
	_, macAddr, found := strings.Cut(t[1], "-")
	if found {
		dp.macAddress = macAddr
	}

	data, err := strconv.ParseFloat(string(m.Payload()), 64)
	if err != nil {
		return dp, err
	}

	dp.metric[metricName] = data

	return dp, nil
}

// Parse topic and fill in values in the data struct
func relayParseVals(dp shelly25Data, metricName string, msg mqtt.Message) (shelly25Data, error) {
	// Basic format checking and filling in values if we have correctly formated Topics
	t := strings.Split(msg.Topic(), "/")

	// Topic lenght should be 5 if not, error out!
	if len(t) != 5 {
		err := errors.New("malformed topic: " + msg.Topic())
		return dp, err
	}
	_, macAddr, found := strings.Cut(t[1], "-")
	if found {
		dp.macAddress = macAddr
	}
	dp.relay = t[3]

	data, err := strconv.ParseFloat(string(msg.Payload()), 64)
	if err != nil {
		return dp, err
	}

	dp.metric[metricName] = data

	return dp, nil
}

// Generic Shelly 2.5 Influx write func
func write25point(dp shelly25Data, metric string) {
	p := influxdb2.NewPointWithMeasurement("Shelly 2.5 Metrics")

	// check if we have a Relay set and add the Tag if we do.
	if dp.relay != "" {
		p.AddTag("Relay", dp.relay)
	}
	p.AddTag("Mac Address", dp.macAddress)
	p.SetTime(time.Now())
	p.AddField(metric, dp.metric[metric])

	database.DbWrite.WritePoint(p)
	health.PointsWritten++
}
