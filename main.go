/*
Based in part from the example the following Blog
https://levelup.gitconnected.com/how-to-use-mqtt-with-go-89c617915774

Bugs and additions added by Joshua Snyder 2022
This is designed to run in Kubernetes. So it features heath, readyness checks and configuration via env vars.
*/
package main

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/cheetahfox/Iot-local-midware/config"
	"github.com/cheetahfox/Iot-local-midware/database"
	mqttcallbacks "github.com/cheetahfox/Iot-local-midware/mqtt_callbacks"
	"github.com/cheetahfox/Iot-local-midware/router"

	mqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/gofiber/fiber/v2"
)

const Version = "0.01"

func main() {
	config := config.Startup()

	iotmw := fiber.New(config.FiberConfig)

	router.SetupRoutes(iotmw)

	go func() {
		iotmw.Listen(":2200")
	}()

	database.ConnectInflux(config)
	mqttcallbacks.SetDefaultCallbacks(&config)
	client := mqtt.NewClient(&config.Options)
	token := client.Connect()
	if token.Wait() && token.Error() != nil {
		panic(token.Error())
	}

	token = client.Subscribe(config.MqttTopic, 1, mqttcallbacks.MessageSubHandler)
	token.Wait()
	fmt.Printf("Subscribed to topic %s\n", config.MqttTopic)

	// Listen for Sigint or SigTerm and exit if you get them.
	sigs := make(chan os.Signal, 1)
	done := make(chan bool, 1)

	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		sig := <-sigs
		fmt.Println()
		fmt.Println(sig)
		done <- true
	}()

	<-done
	fmt.Println("Shutdown Started")
	iotmw.Shutdown()
	database.DisconnectInflux()
	client.Disconnect(100)
}
