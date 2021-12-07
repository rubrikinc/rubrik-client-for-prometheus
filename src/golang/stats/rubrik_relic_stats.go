package stats

import (
	"log"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/rubrikinc/rubrik-sdk-for-go/rubrikcdm"
)

var (
	// storage stats
	rubrikRelicLocalStorage = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "rubrik_relic_local_storage_bytes",
			Help: "Total storage used on local Rubrik cluster by relic objects",
		},
		[]string{
			"ClusterName",
			"ObjectName",
			"ObjectId",
		},
	)

	rubrikRelicArchiveStorage = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "rubrik_relic_archive_storage_bytes",
			Help: "Total storage used in archive locations by relic objects",
		},
		[]string{
			"ClusterName",
			"ObjectName",
			"ObjectId",
		},
	)
)

func init() {
	// Register relic storage stats
	prometheus.MustRegister(rubrikRelicLocalStorage)
	prometheus.MustRegister(rubrikRelicArchiveStorage)
}

// GetRelicStorageStats ...
func GetRelicStorageStats(rubrik *rubrikcdm.Credentials, clusterName string) {
	relicRequest, err := rubrik.Get("v1", "/unmanaged_object?unmanaged_status=Relic", 60)
	if err != nil {
		log.Println("Error from stats.GetRelicStorageStats: ", err)
		return
	}
	relicData := relicRequest.(map[string]interface{})["data"].([]interface{})

	for v := range relicData {
		localStorageBytes, archiveStorageBytes := 0.0, 0.0
		objectName, objectId := "null", "null"

		localStorageBytes = relicData[v].(map[string]interface{})["localStorage"].(float64)
		archiveStorageBytes = relicData[v].(map[string]interface{})["archiveStorage"].(float64)
		objectName = relicData[v].(map[string]interface{})["name"].(string)
		objectId = relicData[v].(map[string]interface{})["id"].(string)

		rubrikRelicLocalStorage.WithLabelValues(
			clusterName,
			objectName,
			objectId,
		).Set(localStorageBytes)

		rubrikRelicArchiveStorage.WithLabelValues(
			clusterName,
			objectName,
			objectId,
		).Set(archiveStorageBytes)
	}
}
