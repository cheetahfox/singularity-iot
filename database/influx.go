package database

import (
	"context"
	"fmt"
	"log"
	"net"
	"os"
	"time"

	"github.com/cheetahfox/singularity-iot/config"
	"github.com/cheetahfox/singularity-iot/health"
	influxdb2 "github.com/influxdata/influxdb-client-go/v2"
	"github.com/influxdata/influxdb-client-go/v2/api"
	"github.com/influxdata/influxdb-client-go/v2/domain"
)

var DbWrite api.WriteAPI
var dbclient influxdb2.Client

func ConnectInflux(config config.Configuration) {
	dbclient = influxdb2.NewClient(config.InfluxdbServer, config.Token)
	dbhealth, err := dbclient.Health(context.Background())
	if (err != nil) && dbhealth.Status == domain.HealthCheckStatusFail {
		log.Panic(err)
	}
	DbWrite = dbclient.WriteAPI(config.Org, config.Bucket)
	errorsCh := DbWrite.Errors()

	// Catch any write errors
	go func() {
		var errorCount int
		for err := range errorsCh {
			fmt.Printf("Influx write error: %s\n", err.Error())
			errorCount++

			/*
				check if the Influxdb database is healthy after seeing a error
				I don't really like this logic I need to come up with something
				simpler and more reliable.
			*/
			if !DbHealthCheck(time.Duration(errorCount) * time.Second) {
				health.InfluxReady = false
				fmt.Println("unhealthy Influxdb")
				if errorCount > config.InfluxMaxError {
					fmt.Println("Maximum Influx error count reached!")
					DisconnectInflux()
					time.Sleep(time.Duration(5) * time.Second)
					os.Exit(1)
				}
			} else {
				// Reset error count if the database is healthy
				errorCount = 0
				health.InfluxReady = true
			}
		}
	}()
	fmt.Printf("Connected to Influxdb %s\n", config.InfluxdbServer)
	health.InfluxReady = true
}

// Check health status after sleeping for a duration
func DbHealthCheck(sleepTime time.Duration) bool {
	time.Sleep(sleepTime)
	dbhealth, err := dbclient.Health(context.Background())
	if (err != nil) || dbhealth.Status == domain.HealthCheckStatusFail {
		return false
	}
	return true
}

// Check current DNS server resolution and sleep up to sd (seconds delay) in 15 second intervals
func checkDns(host string, sd int) bool {

	return true
}

func lookupHost(host string) bool {
	_, err := net.LookupIP(host)
	if err != nil {
		return false
	}

	return true
}

func DisconnectInflux() {
	health.InfluxReady = false
	dbclient.Close()
}
