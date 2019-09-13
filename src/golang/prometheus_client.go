/*
Rubrik Prometheus Client

Requirements:
	Go 1.x (tested with 1.11)
	Rubrik SDK for Go (go get github.com/rubrikinc/rubrik-sdk-for-go)
	Prometheus Client for Go (github.com/prometheus/client_golang)
	Rubrik CDM 3.0+
	Environment variables for rubrik_cdm_node_ip (IP of Rubrik node), rubrik_cdm_username (Rubrik username), rubrik_cdm_password (Rubrik password)
*/

package main

import (
	"log"
	"fmt"
	"net/http"
	"time"
	"github.com/rubrikinc/rubrik-sdk-for-go/rubrikcdm"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

var (
	rubrikTotalStorage = prometheus.NewGauge(prometheus.GaugeOpts{
		Name: "rubrik_total_storage_bytes",
		Help: "Total storage in Rubrik cluster.",
	})
	rubrikUsedStorage = prometheus.NewGauge(prometheus.GaugeOpts{
		Name: "rubrik_used_storage_bytes",
		Help: "Used storage in Rubrik cluster.",
	})
	rubrikAvailableSpace = prometheus.NewGauge(prometheus.GaugeOpts{
		Name: "rubrik_available_storage_bytes",
		Help: "Available storage in Rubrik cluster.",
	})
	rubrikSnapshotStorage = prometheus.NewGauge(prometheus.GaugeOpts{
		Name: "rubrik_snapshot_storage_bytes",
		Help: "Snapshot storage in Rubrik cluster.",
	})
	rubrikLivemountStorage = prometheus.NewGauge(prometheus.GaugeOpts{
		Name: "rubrik_livemount_storage_bytes",
		Help: "Live Mount storage in Rubrik cluster.",
	})
	rubrikMiscStorage = prometheus.NewGauge(prometheus.GaugeOpts{
		Name: "rubrik_misc_storage_bytes",
		Help: "Miscellaneous storage in Rubrik cluster.",
	})
)

func init() {
	// Metrics have to be registered to be exposed:
	prometheus.MustRegister(rubrikTotalStorage)
	prometheus.MustRegister(rubrikUsedStorage)
	prometheus.MustRegister(rubrikAvailableSpace)
	prometheus.MustRegister(rubrikSnapshotStorage)
	prometheus.MustRegister(rubrikLivemountStorage)
	prometheus.MustRegister(rubrikMiscStorage)
}

func main() {
	//start := time.Now()

	rubrik, err := rubrikcdm.ConnectEnv()
	if err != nil {
		log.Fatal(err)
	}

	go func() {
		for {
			fmt.Println("Tick...")
			storageStats,err := rubrik.Get("internal","/stats/system_storage")
			if err != nil {
				log.Fatal(err)
			}
			if total, ok := storageStats.(map[string]interface{})["total"].(float64); ok {
				rubrikTotalStorage.Set(total)
			}
			if used, ok := storageStats.(map[string]interface{})["used"].(float64); ok {
				rubrikUsedStorage.Set(used)
			}
			if avail, ok := storageStats.(map[string]interface{})["available"].(float64); ok {
				rubrikAvailableSpace.Set(avail)
			}
			if snapshot, ok := storageStats.(map[string]interface{})["snapshot"].(float64); ok {
				rubrikSnapshotStorage.Set(snapshot)
			}
			if livemount, ok := storageStats.(map[string]interface{})["liveMount"].(float64); ok {
				rubrikLivemountStorage.Set(livemount)
			}
			if misc, ok := storageStats.(map[string]interface{})["miscellaneous"].(float64); ok {
				rubrikMiscStorage.Set(misc)
			}
			time.Sleep(time.Duration(60) * time.Second)
		}
	}()

	// The Handler function provides a default handler to expose metrics
	// via an HTTP server. "/metrics" is the usual endpoint for that.
	http.Handle("/metrics", promhttp.Handler())
	log.Fatal(http.ListenAndServe(":8080", nil))
}
