package main

import (
	"fmt"
	"github.com/andygrunwald/go-jira"
	"github.com/luohuahuang/qex/config"
	jiraUtils "github.com/luohuahuang/qex/internal/jira"
	"github.com/luohuahuang/qex/monitor"
	"github.com/luohuahuang/qex/protocol"
	"os"
	"time"
)

func main() {
	client := getJiraClient()

	runId := fmt.Sprintf("%s", time.Now().Format("2006-01-02-15:04:05"))
	write(client, config.MapATSignOff, protocol.OKRTypeATSignOff, runId)
	write(client, config.MapATFoundBug, protocol.OKRTypeATFoundBug, runId)
}

func write(client *jira.Client, queries map[string]string, okrType int, runId string) {
	for k, v := range queries {
		if issues, err := search(client, v); err != nil {
			monitor.SendAlert(err)
		} else {
			for _, issue := range issues {
				okr := protocol.Jira{
					RunId:   runId,
					Product: k,
					JiraId:  issue.Key,
					OKRType: okrType,
				}
				jiraUtils.Process(okr)
			}
		}
	}
}

func getJiraClient() *jira.Client {
	tp := jira.BasicAuthTransport{
		Username: os.Getenv("JIRA_USERNAME"),
		Password: os.Getenv("JIRA_SENSITIVE_TOKEN"),
	}

	client, err := jira.NewClient(tp.Client(), config.JiraServer)
	if err != nil {
		monitor.SendAlert(err)
		os.Exit(1)
	}
	return client
}

func search(client *jira.Client, searchString string) ([]jira.Issue, error) {
	last := 0
	var issues []jira.Issue
	for {
		opt := &jira.SearchOptions{
			MaxResults: 1000, // Max results can go up to 1000
			StartAt:    last,
		}

		chunk, resp, err := client.Issue.Search(searchString, opt)
		if err != nil {
			monitor.SendAlert(err)
			return nil, err
		}

		total := resp.Total
		if issues == nil {
			issues = make([]jira.Issue, 0, total)
		}
		issues = append(issues, chunk...)
		last = resp.StartAt + len(chunk)
		if last >= total {
			return issues, nil
		}
	}
}
