package stats

import (
	"log"
	"github.com/rubrikinc/rubrik-sdk-for-go/rubrikcdm"
	"github.com/prometheus/client_golang/prometheus"
)

var (
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
)

func init() {
	// job stats
	prometheus.MustRegister(rubrik24HSucceededJobs)
	prometheus.MustRegister(rubrik24HFailedJobs)
	prometheus.MustRegister(rubrik24HCancelledJobs)
}

// GetMssqlLiveMounts ...
func Get24HJobStats(rubrik *rubrikcdm.Credentials, clusterName string) {
	reportData,err := rubrik.Get("internal","/report?report_template=ProtectionTasksDetails&report_type=Canned", 60) // get our protection tasks details report
	if err != nil {
		log.Printf("Error from stats.Get24HJobStats: ",err)
		return
	}
	reports := reportData.(map[string]interface{})["data"].([]interface{})
	reportID := reports[0].(map[string]interface{})["id"]
	chartData,err := rubrik.Get("internal","/report/"+reportID.(string)+"/chart?chart_id=chart0", 60) // get our chart for the report
	if err != nil {
		log.Printf("Error from stats.Get24HJobStats: ",err)
		return
	}
	for _, v := range chartData.([]interface{}) {
		dataColumns := v.(map[string]interface{})["dataColumns"]
		for _, w := range dataColumns.([]interface{}) {
			label := w.(map[string]interface{})["label"]
			dataPoints := w.(map[string]interface{})["dataPoints"].([]interface{})
			value := dataPoints[0].(map[string]interface{})["value"].(float64)
			switch label {
			case "Succeeded":
				rubrik24HSucceededJobs.WithLabelValues(clusterName).Set(value)
			case "Failed":
				rubrik24HFailedJobs.WithLabelValues(clusterName).Set(value)
			case "Canceled":
				rubrik24HCancelledJobs.WithLabelValues(clusterName).Set(value)
			}
		}
	}
}