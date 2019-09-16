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
	rubrik24HSucceededJobs = prometheus.NewGauge(prometheus.GaugeOpts{
		Name: "rubrik_24h_succeeded_jobs",
		Help: "Last 24 hours succeeded jobs in Rubrik cluster.",
	})
	rubrik24HFailedJobs = prometheus.NewGauge(prometheus.GaugeOpts{
		Name: "rubrik_24h_failed_jobs",
		Help: "Last 24 hours failed jobs in Rubrik cluster.",
	})
	rubrik24HCancelledJobs = prometheus.NewGauge(prometheus.GaugeOpts{
		Name: "rubrik_24h_cancelled_jobs",
		Help: "Last 24 hours cancelled jobs in Rubrik cluster.",
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
	prometheus.MustRegister(rubrik24HSucceededJobs)
	prometheus.MustRegister(rubrik24HFailedJobs)
	prometheus.MustRegister(rubrik24HCancelledJobs)
}

func main() {
	rubrik, err := rubrikcdm.ConnectEnv()
	if err != nil {
		log.Fatal(err)
	}

	// get storage summary
	go func() {
		for {
			storageStats,err := rubrik.Get("internal","/stats/system_storage")
			if err != nil {
				log.Fatal(err)
			}
			// get total storage stat
			if total, ok := storageStats.(map[string]interface{})["total"].(float64); ok {
				rubrikTotalStorage.Set(total)
			}
			// get used storage stat
			if used, ok := storageStats.(map[string]interface{})["used"].(float64); ok {
				rubrikUsedStorage.Set(used)
			}
			// get available storage stat
			if avail, ok := storageStats.(map[string]interface{})["available"].(float64); ok {
				rubrikAvailableSpace.Set(avail)
			}
			// get snapshot storage stat
			if snapshot, ok := storageStats.(map[string]interface{})["snapshot"].(float64); ok {
				rubrikSnapshotStorage.Set(snapshot)
			}
			// get live mount storage stat
			if livemount, ok := storageStats.(map[string]interface{})["liveMount"].(float64); ok {
				rubrikLivemountStorage.Set(livemount)
			}
			// get misc storage stat
			if misc, ok := storageStats.(map[string]interface{})["miscellaneous"].(float64); ok {
				rubrikMiscStorage.Set(misc)
			}
			time.Sleep(time.Duration(60) * time.Second)
		}
	}()

	/*
	// get node stats
	go func() {
		for {
			nodes,err := rubrik.Get("internal","/node")
			if err != nil {
				log.Fatal(err)
			}
			for _, v := range nodes.(map[string]interface{})["data"].([]interface{}) {
				thisNode := (v.(interface{}).(map[string]interface{})["id"])
				nodeDetail,err := rubrik.Get("internal","/node/"+thisNode.(string))
				if err != nil {
					log.Fatal(err)
				}
				fmt.Println(nodeDetail.(map[string]interface{})["status"])
				nodeStats,err := rubrik.Get("internal","/node/"+thisNode.(string)+"/stats")
				if err != nil {
					log.Fatal(err)
				}
				// get cpu stat

				// get network throughput stats

				fmt.Println(nodeStats.(map[string]interface{})["ipAddress"])
			}
			time.Sleep(time.Duration(5) * time.Second)
		}
	}()
	*/

	// get job stats
	go func() {
		for {
			reportData,err := rubrik.Get("internal","/report?report_template=ProtectionTasksDetails&report_type=Canned") // get our protection tasks details report
			if err != nil {
				log.Fatal(err)
			}
			reports := reportData.(map[string]interface{})["data"].([]interface{})
			report_id := reports[0].(map[string]interface{})["id"]
			chartData,err := rubrik.Get("internal","/report/"+report_id.(string)+"/chart?chart_id=chart0") // get our chart for the report
			if err != nil {
				log.Fatal(err)
			}
			for _, v := range chartData.([]interface{}) {
				dataColumns := v.(map[string]interface{})["dataColumns"]
				for _, w := range dataColumns.([]interface{}) {
					label := w.(map[string]interface{})["label"]
					dataPoints := w.(map[string]interface{})["dataPoints"].([]interface{})
					value := dataPoints[0].(map[string]interface{})["value"].(float64)
					switch label {
					case "Succeeded":
						rubrik24HSucceededJobs.Set(value)
					case "Failed":
						rubrik24HFailedJobs.Set(value)
					case "Canceled":
						rubrik24HCancelledJobs.Set(value)
					}
				}
			}
			time.Sleep(time.Duration(1) * time.Hour)
		}
	}()

	// The Handler function provides a default handler to expose metrics
	// via an HTTP server. "/metrics" is the usual endpoint for that.
	http.Handle("/metrics", promhttp.Handler())
	log.Fatal(http.ListenAndServe(":8080", nil))
}
