package main

import (
	"fmt"
	"github.com/andygrunwald/go-jira"
	"github.com/luohuahuang/qex/config"
	"github.com/luohuahuang/qex/internal/influx"
	jiraUtils "github.com/luohuahuang/qex/internal/jira"
	"github.com/luohuahuang/qex/pkg/mattermost"
	"github.com/luohuahuang/qex/protocol"
	"time"
)

func main() {
	client := jiraUtils.GetJiraClient()

	runId := fmt.Sprintf("%s", time.Now().Format("2006-01-02-15:04:05"))
	write(client, config.MapATSignOff, protocol.OKRTypeATSignOff, runId)
	write(client, config.MapATFoundBug, protocol.OKRTypeATFoundBug, runId)
}

func write(client *jira.Client, queries map[string]string, okrType int, runId string) {
	for k, v := range queries {
		if issues, err := jiraUtils.Search(client, v); err != nil {
			mattermost.SendAlert(err, config.MatterMostMonitor)
		} else {
			for _, issue := range issues {
				okr := protocol.Jira{
					RunId:   runId,
					Product: k,
					JiraId:  issue.Key,
					OKRType: okrType,
				}
				influx_utils.ProcessJiraOKR(okr)
			}
		}
	}
}
