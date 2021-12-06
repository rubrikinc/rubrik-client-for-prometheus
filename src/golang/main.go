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
	"log"
	"net/http"
	"os"
	"time"

	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/rubrikinc/rubrik-client-for-prometheus/src/golang/jobs"
	"github.com/rubrikinc/rubrik-client-for-prometheus/src/golang/livemount"
	"github.com/rubrikinc/rubrik-client-for-prometheus/src/golang/objectprotection"
	"github.com/rubrikinc/rubrik-client-for-prometheus/src/golang/stats"
	"github.com/rubrikinc/rubrik-sdk-for-go/rubrikcdm"
)

func main() {
	// set our Prometheus variables
	httpPortEnv, _ := os.LookupEnv("RUBRIK_PROMETHEUS_PORT")
	var httpPort string
	if httpPortEnv == "" {
		httpPort = "8080"
	} else {
		httpPort = httpPortEnv
	}
	rubrik, err := rubrikcdm.ConnectEnv()
	if err != nil {
		log.Printf("Error from main.go:")
		log.Fatal(err)
	}
	clusterDetails, err := rubrik.Get("v1", "/cluster/me", 60)
	if err != nil {
		log.Printf("Error from main.go:")
		log.Fatal(err)
	}
	clusterName := clusterDetails.(map[string]interface{})["name"]
	log.Printf("Cluster name: " + clusterName.(string))

	// get storage summary
	go func() {
		for {
			stats.GetStorageSummaryStats(rubrik, clusterName.(string))
			stats.GetRunwayRemaining(rubrik, clusterName.(string))
			time.Sleep(time.Duration(1) * time.Minute)
		}
	}()

	// get node stats
	go func() {
		for {
			stats.GetNodeStats(rubrik, clusterName.(string))
			time.Sleep(time.Duration(1) * time.Minute)
		}
	}()

	// get job stats
	go func() {
		for {
			stats.Get24HJobStats(rubrik, clusterName.(string))
			time.Sleep(time.Duration(1) * time.Hour)
		}
	}()

	// get compliance stats
	go func() {
		for {
			stats.GetSlaComplianceStats(rubrik, clusterName.(string))
			time.Sleep(time.Duration(1) * time.Hour)
		}
	}()

	// failed job details
	go func() {
		for {
			jobs.GetMssqlFailedJobs(rubrik, clusterName.(string))
			jobs.GetVmwareVmFailedJobs(rubrik, clusterName.(string))
			time.Sleep(time.Duration(5) * time.Minute)
		}
	}()

	// SQL DB capacity stats
	go func() {
		for {
			stats.GetMssqlCapacityStats(rubrik, clusterName.(string))
			time.Sleep(time.Duration(1) * time.Hour)
		}
	}()

	// Oracle DB capacity stats
	go func() {
		for {
			stats.GetOracleCapacityStats(rubrik, clusterName.(string))
			time.Sleep(time.Duration(1) * time.Hour)
		}
	}()

	// VMware vSphere VM capacity stats
	go func() {
		for {
			stats.GetVSphereVmCapacityStats(rubrik, clusterName.(string))
			time.Sleep(time.Duration(1) * time.Hour)
		}
	}()

	// Rubrik Snappable slaDomain
	go func() {
		for {
			objectprotection.GetSnappableEffectiveSlaDomain(rubrik, clusterName.(string))
			time.Sleep(time.Duration(1) * time.Hour)
		}
	}()

	// Rubrik slaDomain Summary
	go func() {
		for {
			objectprotection.GetSlaDomainSummary(rubrik, clusterName.(string))
			time.Sleep(time.Duration(1) * time.Hour)
		}
	}()

	// get live mount stats
	go func() {
		for {
			livemount.GetMssqlLiveMountAges(rubrik, clusterName.(string))
			time.Sleep(time.Duration(1) * time.Hour)
		}
	}()

	//Get Relic Storage Stats
	go func() {
		for {
			stats.GetRelicStorageStats(rubrik, clusterName.(string))
			time.Sleep(time.Duration(1) * time.Hour)
		}
	}()

	// The Handler function provides a default handler to expose metrics
	// via an HTTP server. "/metrics" is the usual endpoint for that.
	http.Handle("/metrics", promhttp.Handler())
	log.Printf("Starting on HTTP port " + httpPort)
	log.Fatal(http.ListenAndServe(":"+httpPort, nil))
}
