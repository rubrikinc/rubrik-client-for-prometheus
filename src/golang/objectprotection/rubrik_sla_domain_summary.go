package objectprotection

import (
	"log"
	"strconv"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/rubrikinc/rubrik-sdk-for-go/rubrikcdm"
)

var (
	// Rubrik SLA Domain Summary Information
	slaDomainSummary = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "rubrik_sla_domain_summary",
			Help: "Return summary information for an SLA domain",
		},
		[]string{
			"primaryClusterId",
			"slaDomainName",
			"slaDomainId",
			"maxLocalRetentionLimit",
			"archivalLocationName",
			"replicationTargetname",
			"hourlyFrequency",
			"hourlyRetention",
			"dailyFrequency",
			"dailyRetention",
			"weeklyFrequency",
			"weeklyRetention",
			"monthlyFrequency",
			"monthlyRetention",
			"quarterlyFrequency",
			"quarterlyRetention",
			"yearlyFrequency",
			"yearlyRetention",
		},
	)
)

func init() {
	prometheus.MustRegister(slaDomainSummary)
}

// GetSlaDomainSummary ...
func GetSlaDomainSummary(rubrik *rubrikcdm.Credentials, clusterName string) {
	slaData, err := rubrik.Get("v2", "/sla_domain", 60) // Get our SLAs

	if err != nil {
		log.Println("Error from objectprotection.GetSlaDomainSummary: ", err)
		return
	}

	slaEntities := slaData.(map[string]interface{})["data"].([]interface{})

	for v := range slaEntities {
		thisClusterId, thisSlaDomainName, thisSlaDomainId := "null", "null", "null"
		thisArchivalLocationName, thisReplicationTargetName := "Not Archived", "Not Replicated"
		thisHourlyFrequency, thisHourlyRetention := 0.0, 0.0
		thisDailyFrequency, thisDailyRetention := 0.0, 0.0
		thisWeeklyFrequency, thisWeeklyRetention := 0.0, 0.0
		thisMonthlyFrequency, thisMonthlyRetention := 0.0, 0.0
		thisQuarterlyFrequency, thisQuarterlyRetention := 0.0, 0.0
		thisYearlyFrequency, thisYearlyRetention := 0.0, 0.0

		thisClusterId = slaEntities[v].(map[string]interface{})["primaryClusterId"].(string)
		thisSlaDomainName = slaEntities[v].(map[string]interface{})["name"].(string)
		thisSlaDomainId = slaEntities[v].(map[string]interface{})["id"].(string)
		thisFrequencies := slaEntities[v].(map[string]interface{})["frequencies"]

		var thisMaxLocalRetention = slaEntities[v].(map[string]interface{})["maxLocalRetentionLimit"].(float64)
		if len(slaEntities[v].(map[string]interface{})["archivalSpecs"].([]interface{})) > 0 {
			thisArchivalLocationName = slaEntities[v].(map[string]interface{})["archivalSpecs"].([]interface{})[0].(map[string]interface{})["locationName"].(string)
		}
		if len(slaEntities[v].(map[string]interface{})["replicationSpecs"].([]interface{})) > 0 {
			thisReplicationTargetName = slaEntities[v].(map[string]interface{})["replicationSpecs"].([]interface{})[0].(map[string]interface{})["locationName"].(string)
		}
		if thisFrequencies.(map[string]interface{})["hourly"] != nil {
			thisHourlyFrequency = thisFrequencies.(map[string]interface{})["hourly"].(map[string]interface{})["frequency"].(float64)
			thisHourlyRetention = thisFrequencies.(map[string]interface{})["hourly"].(map[string]interface{})["retention"].(float64)
		}
		if thisFrequencies.(map[string]interface{})["daily"] != nil {
			thisDailyFrequency = thisFrequencies.(map[string]interface{})["daily"].(map[string]interface{})["frequency"].(float64)
			thisDailyRetention = thisFrequencies.(map[string]interface{})["daily"].(map[string]interface{})["retention"].(float64)
		}
		if thisFrequencies.(map[string]interface{})["weekly"] != nil {
			thisWeeklyFrequency = thisFrequencies.(map[string]interface{})["weekly"].(map[string]interface{})["frequency"].(float64)
			thisWeeklyRetention = thisFrequencies.(map[string]interface{})["weekly"].(map[string]interface{})["retention"].(float64)
		}
		if thisFrequencies.(map[string]interface{})["monthly"] != nil {
			thisMonthlyFrequency = thisFrequencies.(map[string]interface{})["monthly"].(map[string]interface{})["frequency"].(float64)
			thisMonthlyRetention = thisFrequencies.(map[string]interface{})["monthly"].(map[string]interface{})["retention"].(float64)
		}
		if thisFrequencies.(map[string]interface{})["quarterly"] != nil {
			thisQuarterlyFrequency = thisFrequencies.(map[string]interface{})["quarterly"].(map[string]interface{})["frequency"].(float64)
			thisQuarterlyRetention = thisFrequencies.(map[string]interface{})["quarterly"].(map[string]interface{})["retention"].(float64)
		}
		if thisFrequencies.(map[string]interface{})["yearly"] != nil {
			thisYearlyFrequency = thisFrequencies.(map[string]interface{})["yearly"].(map[string]interface{})["frequency"].(float64)
			thisYearlyRetention = thisFrequencies.(map[string]interface{})["yearly"].(map[string]interface{})["retention"].(float64)
		}

		slaDomainSummary.WithLabelValues(
			thisClusterId,
			thisSlaDomainName,
			thisSlaDomainId,
			strconv.FormatFloat(thisMaxLocalRetention, 'f', -1, 64),
			thisArchivalLocationName,
			thisReplicationTargetName,
			strconv.FormatFloat(thisHourlyFrequency, 'f', -1, 64),
			strconv.FormatFloat(thisHourlyRetention, 'f', -1, 64),
			strconv.FormatFloat(thisDailyFrequency, 'f', -1, 64),
			strconv.FormatFloat(thisDailyRetention, 'f', -1, 64),
			strconv.FormatFloat(thisWeeklyFrequency, 'f', -1, 64),
			strconv.FormatFloat(thisWeeklyRetention, 'f', -1, 64),
			strconv.FormatFloat(thisMonthlyFrequency, 'f', -1, 64),
			strconv.FormatFloat(thisMonthlyRetention, 'f', -1, 64),
			strconv.FormatFloat(thisQuarterlyFrequency, 'f', -1, 64),
			strconv.FormatFloat(thisQuarterlyRetention, 'f', -1, 64),
			strconv.FormatFloat(thisYearlyFrequency, 'f', -1, 64),
			strconv.FormatFloat(thisYearlyRetention, 'f', -1, 64),
		).Set(0)
	}
}
