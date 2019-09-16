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
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

var (
	// storage stats
	rubrikTotalStorage = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "rubrik_total_storage_bytes",
			Help: "Total storage in Rubrik cluster.",
		},
		[]string{
			"clusterName",
		},
	)
	rubrikUsedStorage = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "rubrik_used_storage_bytes",
			Help: "Used storage in Rubrik cluster.",
		},
		[]string{
			"clusterName",
		},
	)
	rubrikAvailableSpace = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "rubrik_available_storage_bytes",
			Help: "Available storage in Rubrik cluster.",
		},
		[]string{
			"clusterName",
		},
	)
	rubrikSnapshotStorage = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "rubrik_snapshot_storage_bytes",
			Help: "Snapshot storage in Rubrik cluster.",
		},
		[]string{
			"clusterName",
		},
	)
	rubrikLivemountStorage = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "rubrik_livemount_storage_bytes",
			Help: "Live Mount storage in Rubrik cluster.",
		},
		[]string{
			"clusterName",
		},
	)
	rubrikMiscStorage = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "rubrik_misc_storage_bytes",
			Help: "Miscellaneous storage in Rubrik cluster.",
		},
		[]string{
			"clusterName",
		},
	)
	rubrikRunwayRemaining = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "rubrik_runway_remaining",
			Help: "Runway remaining, in days, on Rubrik cluster.",
		},
		[]string{
			"clusterName",
		},
	)
	// node stats
	rubrikNodeStatus = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "rubrik_node_status",
			Help: "Status of node in Rubrik cluster (1 is OK, 0 is anything else).",
		},
		[]string{
			"clusterName",
			"nodeId",
		},
	)
	// job stats
	rubrik24HSucceededJobs = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "rubrik_24h_succeeded_jobs",
			Help: "Last 24 hours succeeded jobs in Rubrik cluster.",
		},
		[]string{
			"clusterName",
		},
	)
	rubrik24HFailedJobs = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "rubrik_24h_failed_jobs",
			Help: "Last 24 hours failed jobs in Rubrik cluster.",
		},
		[]string{
			"clusterName",
		},
	)
	rubrik24HCancelledJobs = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "rubrik_24h_cancelled_jobs",
			Help: "Last 24 hours cancelled jobs in Rubrik cluster.",
		},
		[]string{
			"clusterName",
		},
	)
	// compliance stats
	rubrikSLACompliantCount = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "rubrik_compliant_object_count",
			Help: "Number of SLA compliant objects in Rubrik cluster.",
		},
		[]string{
			"clusterName",
		},
	)
	rubrikSLANonCompliantCount = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "rubrik_non_compliant_object_count",
			Help: "Number of non-SLA compliant objects in Rubrik cluster.",
		},
		[]string{
			"clusterName",
		},
	)
)

func init() {
	// Metrics have to be registered to be exposed:
	// storage stats
	prometheus.MustRegister(rubrikTotalStorage)
	prometheus.MustRegister(rubrikUsedStorage)
	prometheus.MustRegister(rubrikAvailableSpace)
	prometheus.MustRegister(rubrikSnapshotStorage)
	prometheus.MustRegister(rubrikLivemountStorage)
	prometheus.MustRegister(rubrikMiscStorage)
	prometheus.MustRegister(rubrikRunwayRemaining)
	// node stats
	prometheus.MustRegister(rubrikNodeStatus)
	// job stats
	prometheus.MustRegister(rubrik24HSucceededJobs)
	prometheus.MustRegister(rubrik24HFailedJobs)
	prometheus.MustRegister(rubrik24HCancelledJobs)
	// compliance stats
	prometheus.MustRegister(rubrikSLACompliantCount)
	prometheus.MustRegister(rubrikSLANonCompliantCount)
}

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
			storageStats,err := rubrik.Get("internal","/stats/system_storage")
			if err != nil {
				log.Fatal(err)
			}
			// get total storage stat
			if total, ok := storageStats.(map[string]interface{})["total"].(float64); ok {
				rubrikTotalStorage.WithLabelValues(clusterName.(string)).Set(total)
			}
			// get used storage stat
			if used, ok := storageStats.(map[string]interface{})["used"].(float64); ok {
				rubrikUsedStorage.WithLabelValues(clusterName.(string)).Set(used)
			}
			// get available storage stat
			if avail, ok := storageStats.(map[string]interface{})["available"].(float64); ok {
				rubrikAvailableSpace.WithLabelValues(clusterName.(string)).Set(avail)
			}
			// get snapshot storage stat
			if snapshot, ok := storageStats.(map[string]interface{})["snapshot"].(float64); ok {
				rubrikSnapshotStorage.WithLabelValues(clusterName.(string)).Set(snapshot)
			}
			// get live mount storage stat
			if livemount, ok := storageStats.(map[string]interface{})["liveMount"].(float64); ok {
				rubrikLivemountStorage.WithLabelValues(clusterName.(string)).Set(livemount)
			}
			// get misc storage stat
			if misc, ok := storageStats.(map[string]interface{})["miscellaneous"].(float64); ok {
				rubrikMiscStorage.WithLabelValues(clusterName.(string)).Set(misc)
			}
			runwayRemaining,err := rubrik.Get("internal","/stats/runway_remaining")
			if err != nil {
				log.Fatal(err)
			}
			// get runway remaining stat
			if runway, ok := runwayRemaining.(map[string]interface{})["days"].(float64); ok {
				rubrikRunwayRemaining.WithLabelValues(clusterName.(string)).Set(runway)
			}
			time.Sleep(time.Duration(1) * time.Minute)
		}
	}()

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
				thisNodeStatus := nodeDetail.(map[string]interface{})["status"]
				switch thisNodeStatus {
				case "OK":
					rubrikNodeStatus.WithLabelValues(clusterName.(string),thisNode.(string)).Set(1)
				default:
					rubrikNodeStatus.WithLabelValues(clusterName.(string),thisNode.(string)).Set(0)
				}

				nodeStats,err := rubrik.Get("internal","/node/"+thisNode.(string)+"/stats")
				if err != nil {
					log.Fatal(err)
				}
				// get cpu stat

				// get network throughput stats

				fmt.Println(nodeStats.(map[string]interface{})["ipAddress"])
			}
			time.Sleep(time.Duration(1) * time.Minute)
		}
	}()

	// get job stats
	go func() {
		for {
			reportData,err := rubrik.Get("internal","/report?report_template=ProtectionTasksDetails&report_type=Canned") // get our protection tasks details report
			if err != nil {
				log.Fatal(err)
			}
			reports := reportData.(map[string]interface{})["data"].([]interface{})
			reportID := reports[0].(map[string]interface{})["id"]
			chartData,err := rubrik.Get("internal","/report/"+reportID.(string)+"/chart?chart_id=chart0") // get our chart for the report
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
						rubrik24HSucceededJobs.WithLabelValues(clusterName.(string)).Set(value)
					case "Failed":
						rubrik24HFailedJobs.WithLabelValues(clusterName.(string)).Set(value)
					case "Canceled":
						rubrik24HCancelledJobs.WithLabelValues(clusterName.(string)).Set(value)
					}
				}
			}
			time.Sleep(time.Duration(1) * time.Hour)
		}
	}()

	// get compliance stats
	go func() {
		for {
			reportData,err := rubrik.Get("internal","/report?report_template=SlaComplianceSummary&report_type=Canned") // get our sla compliance summary report
			if err != nil {
				log.Fatal(err)
			}
			reports := reportData.(map[string]interface{})["data"].([]interface{})
			reportID := reports[0].(map[string]interface{})["id"]
			chartData,err := rubrik.Get("internal","/report/"+reportID.(string)+"/chart?chart_id=chart0") // get our chart for the report
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
					case "InCompliance":
						rubrikSLACompliantCount.WithLabelValues(clusterName.(string)).Set(value)
					case "NonCompliance":
						rubrikSLANonCompliantCount.WithLabelValues(clusterName.(string)).Set(value)
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
