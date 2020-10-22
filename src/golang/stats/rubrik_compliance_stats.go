package stats

import (
	"log"
	"github.com/rubrikinc/rubrik-sdk-for-go/rubrikcdm"
	"github.com/prometheus/client_golang/prometheus"
)

var (
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
	// compliance stats
	prometheus.MustRegister(rubrikSLACompliantCount)
	prometheus.MustRegister(rubrikSLANonCompliantCount)
}

// Get ...
func GetSlaComplianceStats(rubrik *rubrikcdm.Credentials, clusterName string) {
	reportData,err := rubrik.Get("internal","/report?report_template=SlaComplianceSummary&report_type=Canned", 60) // get our sla compliance summary report
	if err != nil {
		log.Printf("Error from stats.GetSlaComplianceStats: ",err)
	}
	reports := reportData.(map[string]interface{})["data"].([]interface{})
	reportID := reports[0].(map[string]interface{})["id"]
	chartData,err := rubrik.Get("internal","/report/"+reportID.(string)+"/chart?chart_id=chart0") // get our chart for the report
	if err != nil {
		log.Printf("Error from stats.GetSlaComplianceStats: ",err)
		return
	}
	for _, v := range chartData.([]interface{}) {
		dataColumns := v.(map[string]interface{})["dataColumns"]
		for _, w := range dataColumns.([]interface{}) {
			label := w.(map[string]interface{})["label"]
			dataPoints := w.(map[string]interface{})["dataPoints"].([]interface{})
			value := dataPoints[0].(map[string]interface{})["value"].(float64)
			switch label {
			case "InCompliance":
				rubrikSLACompliantCount.WithLabelValues(clusterName).Set(value)
			case "NonCompliance":
				rubrikSLANonCompliantCount.WithLabelValues(clusterName).Set(value)
			}
		}
	}
}