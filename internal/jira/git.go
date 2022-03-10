package jira_utils

import (
	"github.com/luohuahuang/qex/config"
	influxUtils "github.com/luohuahuang/qex/internal/influx"
	"github.com/luohuahuang/qex/monitor"
	"github.com/luohuahuang/qex/protocol"
)

func Process(jira protocol.Jira) {
	if err := influxUtils.WriteToInflux(config.JiraOKRMeasurement, map[string]string{
		"run_id": jira.RunId,
	}, map[string]interface{}{
		"product": jira.Product,
		"jira_id": jira.JiraId,
		"type":    jira.OKRType,
	}); err != nil {
		monitor.SendAlert(err)
	}
}
