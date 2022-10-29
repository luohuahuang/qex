package influx_utils

import (
	"context"
	influx "github.com/influxdata/influxdb-client-go/v2"
	"github.com/luohuahuang/qex/config"
	"github.com/luohuahuang/qex/pkg/mattermost"
	"github.com/luohuahuang/qex/protocol"
	"strings"
	"time"
)

func WriteToInflux(measurement string, tags map[string]string, fields map[string]interface{}) error {
	ctx := context.Background()
	client := influx.NewClient(config.InfluxUrl, config.InfluxToken)
	client.Options().SetHTTPRequestTimeout(10)
	defer client.Close()
	writeClient := client.WriteAPIBlocking(config.InfluxOrg, config.InfluxBucket)
	p := influx.NewPoint(measurement,
		tags,
		fields,
		time.Now())
	err := writeClient.WritePoint(ctx, p)
	if err != nil {
		mattermost.SendAlert(err, config.MatterMostMonitor)
		return err
	}
	return nil
}

func ProcessGitMaintainer(git protocol.Git) {
	if err := WriteToInflux(config.GitMeasurement, map[string]string{
		"run_id": git.RunId,
	}, map[string]interface{}{
		"commit_id":  git.CommitId,
		"case":       git.Case,
		"maintainer": git.Maintainer,
		"product":    git.Product,
	}); err != nil {
		mattermost.SendAlert(err, config.MatterMostMonitor)
	}
}

func ProcessJiraOKR(jira protocol.Jira) {
	if err := WriteToInflux(config.JiraOKRMeasurement, map[string]string{
		"run_id": jira.RunId,
	}, map[string]interface{}{
		"product": jira.Product,
		"jira_id": jira.JiraId,
		"type":    jira.OKRType,
	}); err != nil {
		mattermost.SendAlert(err, config.MatterMostMonitor)
	}
}

func ProcessGitMR(gitMR protocol.GitMR) {
	if err := WriteToInflux(config.GitMRMeasurement, map[string]string{
		"run_id": gitMR.RunId,
	}, map[string]interface{}{
		"product": gitMR.Product,
		"mr_id":   gitMR.MrID,
		"author":  gitMR.Author,
		"state":   gitMR.State,
	}); err != nil {
		mattermost.SendAlert(err, config.MatterMostMonitor)
	}
}

func ProcessJenkinsBuild(buildDetails protocol.JenkinsBuildExporter) {
	repoArr := strings.Split(buildDetails.RepoUrl, "/")
	if err := WriteToInflux(config.JenkinsBuildMeasurement, map[string]string{
		"user": buildDetails.User,
		"repo": repoArr[len(repoArr)-1],
	}, map[string]interface{}{
		"jenkins_url": buildDetails.BuildUrl,
		"branch":      buildDetails.Branch,
		"job_name":    buildDetails.JobName,
		"timestamp":   buildDetails.Timestamp,
		"is_test_job": buildDetails.IsTestJob,
		"result":      buildDetails.Result,
		"total":       buildDetails.TestDetails.TotalCount,
		"fail_count":  buildDetails.TestDetails.FailCount,
		"skip_count":  buildDetails.TestDetails.SkipCount,
	}); err != nil {
		mattermost.SendAlert(err, config.MatterMostMonitor)
	}
}
