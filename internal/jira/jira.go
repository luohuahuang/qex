package jira_utils

import (
	"errors"
	"fmt"
	"github.com/andygrunwald/go-jira"
	"github.com/luohuahuang/qex/config"
	"github.com/luohuahuang/qex/pkg/mattermost"
	"github.com/luohuahuang/qex/protocol"
	"log"
	"math"
	"os"
	"time"
)

func GetJiraClient() *jira.Client {
	tp := jira.BasicAuthTransport{
		Username: os.Getenv("JIRA_USERNAME"),
		Password: os.Getenv("JIRA_SENSITIVE_TOKEN"),
	}

	client, err := jira.NewClient(tp.Client(), config.JiraServer)
	if err != nil {
		mattermost.SendAlert(err, config.MatterMostMonitor)
	}
	return client
}

func Search(client *jira.Client, searchString string) ([]jira.Issue, error) {
	last := 0
	var issues []jira.Issue
	for {
		opt := &jira.SearchOptions{
			MaxResults: 1000, // Max results can go up to 1000
			StartAt:    last,
		}

		chunk, resp, err := client.Issue.Search(searchString, opt)
		if err != nil {
			mattermost.SendAlert(err, config.MatterMostMonitor)
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

const (
	commentTemplate = `_*Auto Update from QA CI/CD Pipeline*_ 
%s 
%s 
%s%s `
	jobSummarySubTemplate = `
*Job Summary:* 
- Job URL     : %s 
- Job Status  : %s 
- Build Cause : %s 
- Triggered By: %s 
- Time        : %s
`

	RepoSubTemplate = `
*Repo Details:* 
- Repo URL   : %s 
- Branch     : %s 
- SHA1       : %s 
- Bugs linked: %s
`

	TestDetailsTemplate = `
*Test Details* 
- Pass Rate    : %s%% ({*}Total{*}: %d, {*}{color:#00875a}Pass{color}{*}: %d, {*}{color:#de350b}Fail{color}{*}: %d, {*}{color:#ffab00}Skip{color}{*}: %d)`

	TestedRepoTemplate = `
- Tested Repo  : %s 
- Tested Branch: %s 
- Tested SHA1  : %s `

	ColorSuccess  = `{*}{color:#00875a}SUCCESS{color}{*}`
	ColorFail     = `{*}{color:#de350b}FAILURE{color}{*}`
	ColorUnstable = `{*}{color:#ffab00}UNSTABLE{color}{*}`
)

func AddComment(client *jira.Client, issueId string, buildExporter protocol.JenkinsBuildExporter) error {
	var result string
	switch buildExporter.Result {
	case "SUCCESS":
		result = ColorSuccess
	case "UNSTABLE":
		result = ColorUnstable
	default:
		result = ColorFail
	}

	user := buildExporter.User
	if user == "" {
		user = "N/A. Auto Triggered"
	}

	tm := time.Unix(buildExporter.Timestamp/1000, 0)
	jobSummary := fmt.Sprintf(jobSummarySubTemplate,
		buildExporter.BuildUrl,
		result,
		buildExporter.BuildCause,
		user,
		tm.Format("2006-01-02 15:04:05"))

	var bugs string
	for k, _ := range buildExporter.Bugs {
		bugs = fmt.Sprintf("%s %s", k, bugs)
	}
	repoSummary := fmt.Sprintf(RepoSubTemplate,
		buildExporter.RepoUrl,
		buildExporter.Branch,
		buildExporter.Sha1,
		bugs)

	var testDetails string
	var testedRepoDetails string

	if buildExporter.IsTestJob {
		passCount := buildExporter.TestDetails.TotalCount - buildExporter.TestDetails.FailCount - buildExporter.TestDetails.SkipCount

		var passRate string
		if buildExporter.TestDetails.TotalCount > 0 {
			passRate = fmt.Sprintf("%.2f", math.Floor(float64(passCount)/float64(buildExporter.TestDetails.TotalCount)*100))
		} else {
			passRate = "0"
		}
		testDetails = fmt.Sprintf(TestDetailsTemplate,
			passRate,
			buildExporter.TestDetails.TotalCount,
			passCount,
			buildExporter.TestDetails.FailCount,
			buildExporter.TestDetails.SkipCount,
		)

		if buildExporter.TestDetails.TestedRepo != "" {
			testedRepoDetails = fmt.Sprintf(TestedRepoTemplate,
				buildExporter.TestDetails.TestedRepo,
				buildExporter.TestDetails.TestedBranch,
				buildExporter.TestDetails.TestedSha1)
		} else {
			testedRepoDetails = ""
		}
	} else {
		testDetails = "This is a service build deploy job. No test report."
	}

	comment := &jira.Comment{
		Body: fmt.Sprintf(commentTemplate, jobSummary, repoSummary, testDetails, testedRepoDetails),
	}

	_, resp, err := client.Issue.AddComment(issueId, comment)
	if err != nil {
		mattermost.SendAlert(errors.New(fmt.Sprintf("%s: %s", issueId, err.Error())), config.MatterMostMonitor)
		return err
	} else {
		mattermost.SendMsg(buildExporter, issueId)
	}
	log.Println(resp)
	return nil
}
