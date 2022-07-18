package database

import (
	"context"
	"fmt"
	"log"

	"github.com/cheetahfox/Iot-local-midware/config"
	influxdb2 "github.com/influxdata/influxdb-client-go/v2"
	"github.com/influxdata/influxdb-client-go/v2/api"
	"github.com/influxdata/influxdb-client-go/v2/domain"
)

var DbWrite api.WriteAPI
var dbclient influxdb2.Client

func ConnectInflux(config config.Configuration) {
	dbclient = influxdb2.NewClient(config.InfluxdbServer, config.Token)
	health, err := dbclient.Health(context.Background())
	if (err != nil) && health.Status == domain.HealthCheckStatusPass {
		log.Panic(err)
	}
	DbWrite := dbclient.WriteAPI(config.Org, config.Bucket)
	errorsCh := DbWrite.Errors()
	// Catch any write errors
	go func() {
		for err := range errorsCh {
			fmt.Printf("write error: %s\n", err.Error())
		}
	}()
}

func DisconnectInflux() {
	dbclient.Close()
}
