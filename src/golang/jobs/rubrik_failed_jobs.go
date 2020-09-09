package jobs

import (
	"log"
	"strconv"
	"strings"
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
	clusterVersion,err := rubrik.ClusterVersion()
	if err != nil {
		log.Println("Error from jobs.GetMssqlFailedJobs: ",err)
	}
	clusterMajorVersion,err := strconv.ParseInt(strings.Split(clusterVersion,".")[0], 10, 64)
	if err != nil {
		log.Println("Error from jobs.GetMssqlFailedJobs: ",err)
	}
	clusterMinorVersion,err := strconv.ParseInt(strings.Split(clusterVersion,".")[1], 10, 64)
	if err != nil {
		log.Println("Error from jobs.GetMssqlFailedJobs: ",err)
	}
	if (clusterMajorVersion == 5 && clusterMinorVersion < 2) || clusterMajorVersion < 5 { // cluster version is older than 5.1
		eventData,err := rubrik.Get("internal","/event_series?status=Failure&event_type=Backup&object_type=Mssql")
		if err != nil {
			log.Println("Error from jobs.GetMssqlFailedJobs: ",err)
		}
		for _, v := range eventData.(map[string]interface{})["data"].([]interface{}) {
			thisEventSeriesID := v.(map[string]interface{})["eventSeriesId"]
			eventSeriesData,err := rubrik.Get("internal","/event_series/"+thisEventSeriesID.(string))
			if err != nil {
				log.Println("Error from jobs.GetMssqlFailedJobs: ",err)
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
				var thisStartTime string
				if v.(map[string]interface{})["startTime"] == nil {
					thisStartTime = "null"
				} else {
					thisStartTime = v.(map[string]interface{})["startTime"].(string)
				}
				var thisEndTime string
				if v.(map[string]interface{})["endTime"] == nil {
					thisEndTime = "null"
				} else {
					thisEndTime = v.(map[string]interface{})["endTime"].(string)
				}
				var thisLogicalSize string
				if v.(map[string]interface{})["objectLogicalSize"] == nil {
					thisLogicalSize = "null"
				} else {
					thisLogicalSize = strconv.FormatFloat(v.(map[string]interface{})["objectLogicalSize"].(float64), 'f', -1, 64)
				}
				var thisDuration string
				if v.(map[string]interface{})["duration"] == nil {
					thisDuration = "null"
				} else {
					thisDuration = v.(map[string]interface{})["duration"].(string)
				}
				thisEventDate := v.(map[string]interface{})["eventDate"]
				rubrikMssqlFailedJob.WithLabelValues(
					clusterName,
					thisObjectName.(string),
					thisObjectID.(string),
					thisLocation.(string),
					thisStartTime,
					thisEndTime,
					thisLogicalSize,
					thisDuration,
					thisEventDate.(string)).Set(1)
			}
		}
	} else { // cluster version is 5.2 or newer
		eventData,err := rubrik.Get("v1","/event/latest?event_status=Failure&event_type=Backup&object_type=Mssql")
		if err != nil {
			log.Println("Error from jobs.GetMssqlFailedJobs: ",err)
		}
		for _, v := range eventData.(map[string]interface{})["data"].([]interface{}) {
			thisEventSeriesID := v.(map[string]interface{})["latestEvent"].(map[string]interface{})["eventSeriesId"]
			eventSeriesData,err := rubrik.Get("v1","/event_series/"+thisEventSeriesID.(string))
			if err != nil {
				log.Println("Error from jobs.GetMssqlFailedJobs: ",err)
			}
			hasFailedEvent := false
			for _, w := range eventSeriesData.(map[string]interface{})["eventDetailList"].([]interface{}) {
				thisEventStatus := w.(map[string]interface{})["status"]
				if thisEventStatus == "Failure" {
					hasFailedEvent = true
				}
				if hasFailedEvent == true {
					thisObjectName := w.(map[string]interface{})["objectName"]
					thisObjectID := w.(map[string]interface{})["objectId"]
					thisLocation := w.(map[string]interface{})["location"]
					var thisStartTime string
					if w.(map[string]interface{})["startTime"] == nil {
						thisStartTime = "null"
					} else {
						thisStartTime = w.(map[string]interface{})["startTime"].(string)
					}
					var thisEndTime string
					if w.(map[string]interface{})["endTime"] == nil {
						thisEndTime = "null"
					} else {
						thisEndTime = w.(map[string]interface{})["endTime"].(string)
					}
					var thisLogicalSize string
					if w.(map[string]interface{})["logicalSize"] == nil {
						thisLogicalSize = "null"
					} else {
						thisLogicalSize = strconv.FormatFloat(w.(map[string]interface{})["logicalSize"].(float64), 'f', -1, 64)
					}
					var thisDuration string
					if w.(map[string]interface{})["duration"] == nil {
						thisDuration = "null"
					} else {
						thisDuration = w.(map[string]interface{})["duration"].(string)
					}
					rubrikMssqlFailedJob.WithLabelValues(
						clusterName,
						thisObjectName.(string),
						thisObjectID.(string),
						thisLocation.(string),
						thisStartTime,
						thisEndTime,
						thisLogicalSize,
						thisDuration).Set(1)
				}
			}
		}
	}
}