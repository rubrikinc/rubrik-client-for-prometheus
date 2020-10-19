package stats

import (
	"log"
	"github.com/rubrikinc/rubrik-sdk-for-go/rubrikcdm"
	"github.com/prometheus/client_golang/prometheus"
)

var (
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
	rubrikNodeCPU = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "rubrik_node_cpu_ratio",
			Help: "Percentage CPU usage of Rubrik node.",
		},
		[]string{
			"clusterName",
			"nodeId",
		},
	)
	rubrikNodeNetworkReceived = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "rubrik_node_network_received_bytes",
			Help: "Network received byte statistic of Rubrik node.",
		},
		[]string{
			"clusterName",
			"nodeId",
		},
	)
	rubrikNodeNetworkTransmitted = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "rubrik_node_network_transmitted_bytes",
			Help: "Network transmitted byte statistic of Rubrik node.",
		},
		[]string{
			"clusterName",
			"nodeId",
		},
	)
)

func init() {
	// node stats
	prometheus.MustRegister(rubrikNodeStatus)
	prometheus.MustRegister(rubrikNodeCPU)
	prometheus.MustRegister(rubrikNodeNetworkReceived)
	prometheus.MustRegister(rubrikNodeNetworkTransmitted)
}

// GetNodeStats ...
func GetNodeStats(rubrik *rubrikcdm.Credentials, clusterName string) {
	nodes,err := rubrik.Get("internal","/node", 60)
	if err != nil {
		log.Printf("Error from stats.GetNodeStats: ",err)
		return
	}
	for _, v := range nodes.(map[string]interface{})["data"].([]interface{}) {
		thisNode := (v.(interface{}).(map[string]interface{})["id"])
		nodeDetail,err := rubrik.Get("internal","/node/"+thisNode.(string), 60)
		if err != nil {
			log.Printf("Error from stats.GetNodeStats: ",err)
			return
		}
		thisNodeStatus := nodeDetail.(map[string]interface{})["status"]
		switch thisNodeStatus {
		case "OK":
			rubrikNodeStatus.WithLabelValues(clusterName,thisNode.(string)).Set(1)
		default:
			rubrikNodeStatus.WithLabelValues(clusterName,thisNode.(string)).Set(0)
		}

		nodeStats,err := rubrik.Get("internal","/node/"+thisNode.(string)+"/stats?range=-6min", 60)
		if err != nil {
			log.Printf("Error from stats.GetNodeStats: ",err)
			return
		}
		// get cpu stat
		cpuData := nodeStats.(map[string]interface{})["cpuStat"].([]interface{})
		thisCPUStat := cpuData[len(cpuData) - 1].(map[string]interface{})["stat"].(float64) / 100
		rubrikNodeCPU.WithLabelValues(clusterName,thisNode.(string)).Set(thisCPUStat)
		// get network throughput stats
		networkData := nodeStats.(map[string]interface{})["networkStat"]
		byteRxData := networkData.(map[string]interface{})["bytesReceived"].([]interface{})
		thisRxStat := byteRxData[len(byteRxData) - 1].(map[string]interface{})["stat"].(float64)
		rubrikNodeNetworkReceived.WithLabelValues(clusterName,thisNode.(string)).Set(thisRxStat)
		byteTxData := networkData.(map[string]interface{})["bytesTransmitted"].([]interface{})
		thisTxStat := byteTxData[len(byteTxData) - 1].(map[string]interface{})["stat"].(float64)
		rubrikNodeNetworkTransmitted.WithLabelValues(clusterName,thisNode.(string)).Set(thisTxStat)
	}
}