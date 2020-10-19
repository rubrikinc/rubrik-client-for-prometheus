package livemount

import (
	"log"
	"time"
	"github.com/rubrikinc/rubrik-sdk-for-go/rubrikcdm"
	"github.com/prometheus/client_golang/prometheus"
)

var (
	// live mount stats
	rubrikMssqlLiveMountAge = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "rubrik_mssql_live_mount_age_seconds",
			Help: "Age of SQL DB live mounts.",
		},
		[]string{
			"clusterName",
			"sourceDatabaseName",
			"sourceDatabaseId",
			"mountedDatabaseName",
		},
	)
)

func init() {
	// live mount stats
	prometheus.MustRegister(rubrikMssqlLiveMountAge)
}

// GetMssqlLiveMountAges ...
func GetMssqlLiveMountAges(rubrik *rubrikcdm.Credentials, clusterName string) {
	mountData,err := rubrik.Get("v1","/mssql/db/mount", 60) // get our mssql live mount summary
	if err != nil {
		log.Printf("Error from livemount.GetMssqlLiveMountAges: ",err)
		return
	}
	for _, v := range mountData.(map[string]interface{})["data"].([]interface{}) {
		thisSourceDatabaseName := v.(map[string]interface{})["sourceDatabaseName"]
		thisSourceDatabaseID := v.(map[string]interface{})["sourceDatabaseId"]
		thisMountedDatabaseName := v.(map[string]interface{})["mountedDatabaseName"]
		thisCreationDate := v.(map[string]interface{})["creationDate"]
		mountTime, _ := time.Parse(time.RFC3339, thisCreationDate.(string))
		age := time.Since(mountTime)
		//fmt.Println(age.Seconds())
		rubrikMssqlLiveMountAge.WithLabelValues(
			clusterName,
			thisSourceDatabaseName.(string),
			thisSourceDatabaseID.(string),
			thisMountedDatabaseName.(string)).Set(age.Seconds())
	}
}