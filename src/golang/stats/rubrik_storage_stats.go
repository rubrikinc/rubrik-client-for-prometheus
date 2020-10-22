package stats

import (
	"log"
	"github.com/rubrikinc/rubrik-sdk-for-go/rubrikcdm"
	"github.com/prometheus/client_golang/prometheus"
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
)

func init() {
	// storage stats
	prometheus.MustRegister(rubrikTotalStorage)
	prometheus.MustRegister(rubrikUsedStorage)
	prometheus.MustRegister(rubrikAvailableSpace)
	prometheus.MustRegister(rubrikSnapshotStorage)
	prometheus.MustRegister(rubrikLivemountStorage)
	prometheus.MustRegister(rubrikMiscStorage)
	prometheus.MustRegister(rubrikRunwayRemaining)
}

// GetStorageSummaryStats ...
func GetStorageSummaryStats(rubrik *rubrikcdm.Credentials, clusterName string) {
	storageStats,err := rubrik.Get("internal","/stats/system_storage", 60)
	if err != nil {
		log.Printf("Error from stats.GetStorageSummaryStats: ",err)
		return
	}
	// get total storage stat
	if total, ok := storageStats.(map[string]interface{})["total"].(float64); ok {
		rubrikTotalStorage.WithLabelValues(clusterName).Set(total)
	}
	// get used storage stat
	if used, ok := storageStats.(map[string]interface{})["used"].(float64); ok {
		rubrikUsedStorage.WithLabelValues(clusterName).Set(used)
	}
	// get available storage stat
	if avail, ok := storageStats.(map[string]interface{})["available"].(float64); ok {
		rubrikAvailableSpace.WithLabelValues(clusterName).Set(avail)
	}
	// get snapshot storage stat
	if snapshot, ok := storageStats.(map[string]interface{})["snapshot"].(float64); ok {
		rubrikSnapshotStorage.WithLabelValues(clusterName).Set(snapshot)
	}
	// get live mount storage stat
	if livemount, ok := storageStats.(map[string]interface{})["liveMount"].(float64); ok {
		rubrikLivemountStorage.WithLabelValues(clusterName).Set(livemount)
	}
	// get misc storage stat
	if misc, ok := storageStats.(map[string]interface{})["miscellaneous"].(float64); ok {
		rubrikMiscStorage.WithLabelValues(clusterName).Set(misc)
	}
}

// GetRunwayRemaining ...q
func GetRunwayRemaining(rubrik *rubrikcdm.Credentials, clusterName string) {
	runwayRemaining,err := rubrik.Get("internal","/stats/runway_remaining", 60)
	if err != nil {
		log.Printf("Error from stats.GetRunwayRemaining: ",err)
		return
	}
	// get runway remaining stat
	if runway, ok := runwayRemaining.(map[string]interface{})["days"].(float64); ok {
		rubrikRunwayRemaining.WithLabelValues(clusterName).Set(runway)
	}
}