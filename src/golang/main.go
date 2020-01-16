/*
Rubrik Prometheus Client

Requirements:
	Go 1.x (tested with 1.11)
	Rubrik SDK for Go (go get github.com/rubrikinc/rubrik-sdk-for-go)
	Prometheus Client for Go (go get github.com/prometheus/client_golang)
	Rubrik CDM 3.0+
	Environment variables for rubrik_cdm_node_ip (IP of Rubrik node), rubrik_cdm_username (Rubrik username), rubrik_cdm_password (Rubrik password)
*/

package main

import (
	"fmt"
	"log"
	"net/http"
	"time"
	"github.com/rubrikinc/rubrik-sdk-for-go/rubrikcdm"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

func main() {
	rubrik, err := rubrikcdm.ConnectEnv()
	if err != nil {
		log.Fatal(err)
	}
	clusterDetails,err := rubrik.Get("v1","/cluster/me")
	if err != nil {
		log.Fatal(err)
	}
	clusterName := clusterDetails.(map[string]interface{})["name"]
	fmt.Println("Cluster name: "+clusterName.(string))

	// get storage summary
	go func() {
		for {
			GetStorageSummaryStats(rubrik, clusterName.(string))
			GetRunwayRemaining(rubrik, clusterName.(string))
			time.Sleep(time.Duration(1) * time.Minute)
		}
	}()

	// get node stats
	go func() {
		for {
			GetNodeStats(rubrik, clusterName.(string))
			time.Sleep(time.Duration(1) * time.Minute)
		}
	}()

	// get job stats
	go func() {
		for {
			Get24HJobStats(rubrik, clusterName.(string))
			time.Sleep(time.Duration(1) * time.Hour)
		}
	}()

	// get compliance stats
	go func() {
		for {
			GetSlaComplianceStats(rubrik, clusterName.(string))
			time.Sleep(time.Duration(1) * time.Hour)
		}
	}()

	// failed job details
	go func() {
		for {
			GetMssqlFailedJobs(rubrik, clusterName.(string))
			time.Sleep(time.Duration(5) * time.Minute)
		}
	}()

	// SQL DB capacity stats
	go func() {
		for {
			GetMssqlCapacityStats(rubrik, clusterName.(string))
			time.Sleep(time.Duration(1) * time.Hour)
		}
	}()

	// get live mount stats
	go func() {
		for {
			GetMssqlLiveMountAges(rubrik, clusterName.(string))
			time.Sleep(time.Duration(1) * time.Hour)
		}
	}()

	// The Handler function provides a default handler to expose metrics
	// via an HTTP server. "/metrics" is the usual endpoint for that.
	http.Handle("/metrics", promhttp.Handler())
	log.Fatal(http.ListenAndServe(":8080", nil))
}
