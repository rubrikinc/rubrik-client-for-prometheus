package jobs

import (
	"log"
	"strconv"
	"github.com/rubrikinc/rubrik-sdk-for-go/rubrikcdm"
	"github.com/prometheus/client_golang/prometheus"
)

var (
	// failed job details
	rubrikMssqlFailedJob = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "rubrik_mssql_failed_job",
			Help: "Information for failed Rubrik MSSQL Backup job.",
		},
		[]string{
			"clusterName",
			"objectName",
			"objectID",
			"location",
			"startTime",
			"endTime",
			"objectLogicalSize",
			"duration",
			"eventDate",
		},
	)
)

func init() {
	// failed job details
	prometheus.MustRegister(rubrikMssqlFailedJob)
}

// GetMssqlFailedJobs ...
func GetMssqlFailedJobs(rubrik *rubrikcdm.Credentials, clusterName string) {
	eventData,err := rubrik.Get("internal","/event_series?status=Failure&event_type=Backup&object_type=Mssql")
	if err != nil {
		log.Fatal(err)
	}

	for _, v := range eventData.(map[string]interface{})["data"].([]interface{}) {
		thisEventSeriesID := v.(map[string]interface{})["eventSeriesId"]
		eventSeriesData,err := rubrik.Get("internal","/event_series/"+thisEventSeriesID.(string))
		if err != nil {
			log.Fatal(err)
		}
		hasFailedEvent := false
		for _, w := range eventSeriesData.(map[string]interface{})["eventDetailList"].([]interface{}) {
			thisEventStatus := w.(map[string]interface{})["status"]
			if thisEventStatus == "Failure" {
				hasFailedEvent = true
			}
		}
		if hasFailedEvent == true {
			thisObjectName := v.(map[string]interface{})["objectInfo"].(map[string]interface{})["objectName"]
			thisObjectID := v.(map[string]interface{})["objectInfo"].(map[string]interface{})["objectId"]
			thisLocation := v.(map[string]interface{})["location"]
			thisStartTime := v.(map[string]interface{})["startTime"]
			if thisStartTime == nil { thisStartTime = "null" }
			thisEndTime := v.(map[string]interface{})["endTime"]
			if thisEndTime == nil { thisEndTime = "null" }
			thisLogicalSize := v.(map[string]interface{})["objectLogicalSize"]
			if thisLogicalSize == nil {
				thisLogicalSize = "null"
			} else {
				thisLogicalSize = strconv.FormatFloat(thisLogicalSize.(float64), 'f', -1, 64)
			}
			thisDuration := v.(map[string]interface{})["duration"]
			if thisDuration == nil { thisDuration = "null" }
			thisEventDate := v.(map[string]interface{})["eventDate"]
			rubrikMssqlFailedJob.WithLabelValues(
				clusterName,
				thisObjectName.(string),
				thisObjectID.(string),
				thisLocation.(string),
				thisStartTime.(string),
				thisEndTime.(string),
				thisLogicalSize.(string),
				thisDuration.(string),
				thisEventDate.(string)).Set(1)
		}
	}
}