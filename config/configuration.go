/*

Yaml device config file format

Default Location: /etc/singularity/devices.yml

---
iotdevices:
  - name:   shellytv
	hwtype: shelly25
	mqid:   "98CDAC38E9F5"
	topic:  "shellies/shellyswitch25-98CDAC38E9F5"
  - name: ...
*/
package config

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strconv"
	"time"

	mqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/gofiber/fiber/v2"
	"gopkg.in/yaml.v2"
)

const DeviceConfig = "/etc/singularity/devices.yml"

// Iotdevices
type Iotdevices struct {
	Name   string `yaml:"name"`
	Hwtype string `yaml:"hwtype"`
	Mqid   string `yaml:"mqid"`
	Topic  string `yaml:"topic"`
}

type Configuration struct {
	Options        mqtt.ClientOptions
	FiberConfig    fiber.Config
	Bucket         string
	InfluxdbServer string
	Org            string
	Token          string
	MqttTopic      string
	Devices        []Iotdevices `yaml:"iotdevices"`
	InfluxMaxError int
}

// Set configuration options from Env values and setup the Fiber options
func Startup() Configuration {
	var conf Configuration

	requiredEnvVars := []string{
		"INFLUX_SERVER", // Influxdb server url including port number
		"INFLUX_TOKEN",  // Influx Token
		"INFLUX_BUCKET", // Influx bucket
		"INFLUX_ORG",    // Influx ord
		"MQTT_BROKER",   // MQTT access Url
		"MQTT_ID",       // MQTT ID of this client
		"MQTT_TOPIC",    // MQTT Topic
		"DB_MAX_ERROR",
	}

	// Check if the Required Enviromental varibles are set exit if they aren't.
	for index := range requiredEnvVars {
		if os.Getenv(requiredEnvVars[index]) == "" {
			log.Fatalf("Missing %s Enviroment var \n", requiredEnvVars[index])
		}
	}

	conf.Options = *mqtt.NewClientOptions()
	conf.Options.AddBroker(os.Getenv("MQTT_BROKER"))
	conf.Options.SetClientID(os.Getenv("MQTT_ID"))
	conf.MqttTopic = os.Getenv("MQTT_TOPIC")
	// Set a default max error and then read the env value.
	conf.InfluxMaxError = 10
	influxerrors, err := strconv.Atoi(os.Getenv("DB_MAX_ERROR"))
	if err != nil {
		conf.InfluxMaxError = influxerrors
	}

	// Fiber Setup
	conf.FiberConfig = fiber.Config{
		CaseSensitive: true,
		StrictRouting: true,
		ServerHeader:  "IoT-gw",
		AppName:       "IoT Midware v0.01",
		ReadTimeout:   (30 * time.Second),
	}

	// Influxdb Settings
	conf.Token = os.Getenv("INFLUX_TOKEN")
	conf.Bucket = os.Getenv("INFLUX_BUCKET")
	conf.Org = os.Getenv("INFLUX_ORG")
	conf.InfluxdbServer = os.Getenv("INFLUX_SERVER")

	// Read and parse the iot device configuration
	configfile, err := ioutil.ReadFile(DeviceConfig)
	if err != nil {
		fmt.Println("Unable to read device config")
		log.Fatal(err)
	}

	err = yaml.Unmarshal(configfile, &conf)
	if err != nil {
		log.Fatal(err)
	}

	for i, _ := range conf.Devices {
		fmt.Println(conf.Devices[i].Name)
	}

	return conf
}
