package objectprotection

import (
	"log"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/rubrikinc/rubrik-sdk-for-go/rubrikcdm"
)

var (
	// Rubrik Snappable SlaDomain Information
	snappableEffectiveSlaDomain = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "rubrik_snappable_effective_sla",
			Help: "Return the slaDomain information for snappables",
		},
		[]string{
			"clusterName",
			"objectName",
			"objectType",
			"objectID",
			"location",
			"slaDomain",
		},
	)
)

func init() {
	// VMware vSphere VM effective sla domain
	prometheus.MustRegister(snappableEffectiveSlaDomain)
}

// GetSnappableEffectiveSlaDomain ...
func GetSnappableEffectiveSlaDomain(rubrik *rubrikcdm.Credentials, clusterName string) {
	reportData, err := rubrik.Get("internal", "/report?report_template=ObjectProtectionSummary&report_type=Canned", 60) // get our object protection summary report
	if err != nil {
		log.Printf("Error from objectprotection.GetSnappableEffectiveSlaDomain: ", err)
		return
	}
	reports := reportData.(map[string]interface{})["data"].([]interface{})
	reportID := reports[0].(map[string]interface{})["id"]
	body := map[string]interface{}{
		"limit": 100,
	}
	for {
		hasMore := true
		tableData, err := rubrik.Post("internal", "/report/"+reportID.(string)+"/table", body, 60) // get our first page of data for the report
		if err != nil {
			log.Printf("Error from objectprotection.GetSnappableEffectiveSlaDomain: ", err)
			return
		}
		dataGrid := tableData.(map[string]interface{})["dataGrid"].([]interface{})
		hasMore = tableData.(map[string]interface{})["hasMore"].(bool)
		cursor := tableData.(map[string]interface{})["cursor"]
		columns := tableData.(map[string]interface{})["columns"].([]interface{})
		for _, v := range dataGrid {
			thisObjectID, thisObjectName, thisObjectType, thisLocation, thisSlaDomain := "null", "null", "null", "null", "null"

			for i := 0; i < len(columns); i++ {
				switch columns[i] {
				case "ObjectId":
					thisObjectID = v.([]interface{})[i].(string)
				case "ObjectLinkingId":
					thisObjectID = v.([]interface{})[i].(string)
				case "ObjectName":
					thisObjectName = v.([]interface{})[i].(string)
				case "ObjectType":
					thisObjectType = v.([]interface{})[i].(string)
				case "Location":
					thisLocation = v.([]interface{})[i].(string)
				case "SlaDomain":
					thisSlaDomain = v.([]interface{})[i].(string)
				}
			}
			snappableEffectiveSlaDomain.WithLabelValues(
				clusterName,
				thisObjectName,
				thisObjectType,
				thisObjectID,
				thisLocation,
				thisSlaDomain).Set(0)
		}
		if !hasMore {
			return
		} else {
			body = map[string]interface{}{
				"limit":  1000,
				"cursor": cursor,
			}
		}
	}
}
