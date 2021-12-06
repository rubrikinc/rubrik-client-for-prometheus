package jobs

import (
	"log"
	"strconv"
	"strings"
	"time"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/rubrikinc/rubrik-sdk-for-go/rubrikcdm"
)

var (
	// Mssql failed job details
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
	// VM failed job details
	rubrikVmwareVmFailedJob = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "rubrik_vmwarevm_failed_job",
			Help: "Information for failed Rubrik VMware VM Backup job.",
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
	prometheus.MustRegister(rubrikVmwareVmFailedJob)
}

// GetMssqlFailedJobs ...
func GetMssqlFailedJobs(rubrik *rubrikcdm.Credentials, clusterName string) {
	clusterVersion, err := rubrik.ClusterVersion()
	if err != nil {
		log.Printf("Error from jobs.GetMssqlFailedJobs: ", err)
		return
	}
	clusterMajorVersion, err := strconv.ParseInt(strings.Split(clusterVersion, ".")[0], 10, 64)
	if err != nil {
		log.Printf("Error from jobs.GetMssqlFailedJobs: ", err)
		return
	}
	clusterMinorVersion, err := strconv.ParseInt(strings.Split(clusterVersion, ".")[1], 10, 64)
	if err != nil {
		log.Printf("Error from jobs.GetMssqlFailedJobs: ", err)
		return
	}
	if (clusterMajorVersion == 5 && clusterMinorVersion < 2) || clusterMajorVersion < 5 { // cluster version is older than 5.1
		eventData, err := rubrik.Get("internal", "/event_series?status=Failure&event_type=Backup&object_type=Mssql", 60)
		if err != nil {
			log.Printf("Error from jobs.GetMssqlFailedJobs: ", err)
			return
		}
		if eventData != nil || eventData.(map[string]interface{})["data"] != nil {
			for _, v := range eventData.(map[string]interface{})["data"].([]interface{}) {
				thisEventSeriesID := v.(map[string]interface{})["eventSeriesId"]
				eventSeriesData, err := rubrik.Get("internal", "/event_series/"+thisEventSeriesID.(string), 60)
				if err != nil {
					log.Printf("Error from jobs.GetMssqlFailedJobs: ", err)
					return
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
		}
	} else { // cluster version is 5.2 or newer
		var yesterday = time.Now().AddDate(0, 0, -1).Format("2006-01-02T15:04:05.000Z")
		eventData, err := rubrik.Get("v1", "/event/latest?limit=9999&event_status=Failure&event_type=Backup&object_type=Mssql&before_date="+yesterday, 60)
		if err != nil {
			log.Printf("Error from jobs.GetMssqlFailedJobs: ", err)
			return
		}
		if eventData != nil || eventData.(map[string]interface{})["data"] != nil {
			for _, v := range eventData.(map[string]interface{})["data"].([]interface{}) {
				thisEventSeriesID := v.(map[string]interface{})["latestEvent"].(map[string]interface{})["eventSeriesId"]
				eventSeriesData, err := rubrik.Get("v1", "/event_series/"+thisEventSeriesID.(string), 60)
				if err != nil {
					log.Printf("Error from jobs.GetMssqlFailedJobs: ", err)
					return
				}
				hasFailedEvent := false
				for _, w := range eventSeriesData.(map[string]interface{})["eventDetailList"].([]interface{}) {
					thisEventStatus := w.(map[string]interface{})["eventStatus"]
					if thisEventStatus == "Failure" {
						hasFailedEvent = true
					}
				}
				if hasFailedEvent == true {
					thisObjectName := eventSeriesData.(map[string]interface{})["objectName"]
					thisObjectID := eventSeriesData.(map[string]interface{})["objectId"]
					thisLocation := eventSeriesData.(map[string]interface{})["location"]
					var thisStartTime string
					if eventSeriesData.(map[string]interface{})["startTime"] == nil {
						thisStartTime = "null"
					} else {
						thisStartTime = eventSeriesData.(map[string]interface{})["startTime"].(string)
					}
					var thisEndTime string
					if eventSeriesData.(map[string]interface{})["endTime"] == nil {
						thisEndTime = "null"
					} else {
						thisEndTime = eventSeriesData.(map[string]interface{})["endTime"].(string)
					}
					var thisLogicalSize string
					if eventSeriesData.(map[string]interface{})["logicalSize"] == nil {
						thisLogicalSize = "null"
					} else {
						thisLogicalSize = strconv.FormatFloat(eventSeriesData.(map[string]interface{})["logicalSize"].(float64), 'f', -1, 64)
					}
					var thisDuration string
					if eventSeriesData.(map[string]interface{})["duration"] == nil {
						thisDuration = "null"
					} else {
						thisDuration = eventSeriesData.(map[string]interface{})["duration"].(string)
					}
					thisEventDate := eventSeriesData.(map[string]interface{})["startTime"]
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
		}
	}
}
// GetVmwareVMFailedJobs ...
func GetVmwareVmFailedJobs(rubrik *rubrikcdm.Credentials, clusterName string) {
	clusterVersion, err := rubrik.ClusterVersion()
	if err != nil {
		log.Printf("Error from jobs.GetVmwareVmFailedJobs: ", err)
		return
	}
	clusterMajorVersion, err := strconv.ParseInt(strings.Split(clusterVersion, ".")[0], 10, 64)
	if err != nil {
		log.Printf("Error from jobs.GetVmwareVmFailedJobs: ", err)
		return
	}
	clusterMinorVersion, err := strconv.ParseInt(strings.Split(clusterVersion, ".")[1], 10, 64)
	if err != nil {
		log.Printf("Error from jobs.GetVmwareVmFailedJobs: ", err)
		return
	}
	if (clusterMajorVersion == 5 && clusterMinorVersion < 2) || clusterMajorVersion < 5 { // cluster version is older than 5.1
		eventData, err := rubrik.Get("internal", "/event_series?status=Failure&event_type=Backup&object_type=VmwareVm", 60)
		if err != nil {
			log.Printf("Error from jobs.GetVmwareVmFailedJobs: ", err)
			return
		}
		if eventData != nil || eventData.(map[string]interface{})["data"] != nil {
			for _, v := range eventData.(map[string]interface{})["data"].([]interface{}) {
				thisEventSeriesID := v.(map[string]interface{})["eventSeriesId"]
				eventSeriesData, err := rubrik.Get("internal", "/event_series/"+thisEventSeriesID.(string), 60)
				if err != nil {
					log.Printf("Error from jobs.GetVmwareVmFailedJobs: ", err)
					return
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
					rubrikVmwareVmFailedJob.WithLabelValues(
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
		}
	} else { // cluster version is 5.2 or newer
		var yesterday = time.Now().AddDate(0, 0, -1).Format("2006-01-02T15:04:05.000Z")
		eventData, err := rubrik.Get("v1", "/event/latest?limit=9999&event_status=Failure&event_type=Backup&object_type=VmwareVm&before_date="+yesterday, 60)
		if err != nil {
			log.Printf("Error from jobs.GetVmwareVmFailedJobs: ", err)
			return
		}
		if eventData != nil || eventData.(map[string]interface{})["data"] != nil {
			for _, v := range eventData.(map[string]interface{})["data"].([]interface{}) {
				thisEventSeriesID := v.(map[string]interface{})["latestEvent"].(map[string]interface{})["eventSeriesId"]
				eventSeriesData, err := rubrik.Get("v1", "/event_series/"+thisEventSeriesID.(string), 60)
				if err != nil {
					log.Printf("Error from jobs.GetVmwareVmFailedJobs: ", err)
					return
				}
				hasFailedEvent := false
				for _, w := range eventSeriesData.(map[string]interface{})["eventDetailList"].([]interface{}) {
					thisEventStatus := w.(map[string]interface{})["eventStatus"]
					if thisEventStatus == "Failure" {
						hasFailedEvent = true
					}
				}
				if hasFailedEvent == true {
					thisObjectName := eventSeriesData.(map[string]interface{})["objectName"]
					thisObjectID := eventSeriesData.(map[string]interface{})["objectId"]
					thisLocation := eventSeriesData.(map[string]interface{})["location"]
					var thisStartTime string
					if eventSeriesData.(map[string]interface{})["startTime"] == nil {
						thisStartTime = "null"
					} else {
						thisStartTime = eventSeriesData.(map[string]interface{})["startTime"].(string)
					}
					var thisEndTime string
					if eventSeriesData.(map[string]interface{})["endTime"] == nil {
						thisEndTime = "null"
					} else {
						thisEndTime = eventSeriesData.(map[string]interface{})["endTime"].(string)
					}
					var thisLogicalSize string
					if eventSeriesData.(map[string]interface{})["logicalSize"] == nil {
						thisLogicalSize = "null"
					} else {
						thisLogicalSize = strconv.FormatFloat(eventSeriesData.(map[string]interface{})["logicalSize"].(float64), 'f', -1, 64)
					}
					var thisDuration string
					if eventSeriesData.(map[string]interface{})["duration"] == nil {
						thisDuration = "null"
					} else {
						thisDuration = eventSeriesData.(map[string]interface{})["duration"].(string)
					}
					thisEventDate := eventSeriesData.(map[string]interface{})["startTime"]
					rubrikVmwareVmFailedJob.WithLabelValues(
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
		}
	}
}
