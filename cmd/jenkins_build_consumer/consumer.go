package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/luohuahuang/qex/config"
	"github.com/luohuahuang/qex/internal/cache"
	gitUtils "github.com/luohuahuang/qex/internal/git"
	"github.com/luohuahuang/qex/internal/influx"
	jiraUtils "github.com/luohuahuang/qex/internal/jira"
	"github.com/luohuahuang/qex/internal/kafka"
	"github.com/luohuahuang/qex/monitor"
	"github.com/luohuahuang/qex/protocol"
	"github.com/xanzy/go-gitlab"
	"io/ioutil"
	"log"
	"net/http"
	"regexp"
	"strconv"
	"strings"
	"time"
)

func ProcessEvent() {
	Consume(kafka.New(config.JenkinsBuildTopic, config.JenkinsZooKeeper))
}

func main() {
	ProcessEvent()
}

func Consume(k *kafka.Consumer) {
	for {
		select {
		case msg := <-k.ConsumerGroup.Messages():
			if msg == nil || msg.Topic != k.Topic {
				continue
			}
			log.Println(fmt.Sprintf("topic: %s; msg: %s", msg.Topic, string(msg.Value)))

			process(string(msg.Value))
			err := k.ConsumerGroup.CommitUpto(msg)
			if err != nil {
				log.Println("error commit zookeeper: ", err.Error())
			}
		}
	}
}

func process(msg string) {

	build, buildDetails := RetrieveBuildDetails(msg)

	buildExporter := protocol.JenkinsBuildExporter{
		BuildUrl:  buildDetails.URL,
		JobName:   build.JobName,
		BuildNum:  build.BuildNum,
		Timestamp: buildDetails.Timestamp,
		Result:    buildDetails.Result,
		IsTestJob: false,
		Bugs:      map[string]string{},
		TestDetails: protocol.TestDetails{
		},
	}

	var isByPipeline bool
	for _, action := range buildDetails.Actions {
		if action.Class == "hudson.tasks.junit.TestResultAction" {
			buildExporter.IsTestJob = true
			buildExporter.TestDetails.TotalCount = action.TotalCount
			buildExporter.TestDetails.FailCount = action.FailCount
			buildExporter.TestDetails.SkipCount = action.SkipCount
		}
		if action.Class == "hudson.plugins.git.util.BuildData" {
			buildExporter.Branch = action.LastBuiltRevision.Branch[0].Name
			buildExporter.Sha1 = action.LastBuiltRevision.Branch[0].Sha1
			buildExporter.RepoUrl = action.RemoteUrls[0]
		}
		if action.Class == "hudson.model.CauseAction" {
			if action.Causes[0].Class == "org.jenkinsci.plugins.workflow.support.steps.build.BuildUpstreamCause" {
				isByPipeline = true
			}
			buildExporter.User = action.Causes[0].UserName
			buildExporter.BuildCause = action.Causes[0].ShortDescription
		}
	}

	redisCli := cache.New(config.CacheServer)
	var key string

	if !buildExporter.IsTestJob {
		var upstreamJob, upstreamBuild string
		if isByPipeline {
			// "Started by upstream project \"my-abc-service-build-deploy\" build number 39"
			temp := strings.Split(buildExporter.BuildCause, `"`)
			if len(temp) != 3 {
				monitor.SendAlert(errors.New("unexpected build cause found: " + buildExporter.BuildCause))
			} else {
				if len(temp) > 1 {
					upstreamJob = temp[1]
				}
			}
			r, _ := regexp.Compile("([0-9]+)$")
			upstreamBuild = r.FindString(buildExporter.BuildCause)
			key = fmt.Sprintf("qex_%s_%s", upstreamJob, upstreamBuild)

			// extract from MR title
			for _, bug := range ExtractBugsFromMRTitle(upstreamJob, upstreamBuild) {
				buildExporter.Bugs[bug] = ""
			}

		} else {
			key = fmt.Sprintf("qex_%s_%s", buildExporter.JobName, buildExporter.BuildNum)
		}

		_ = redisCli.HSet(key, "repo", buildExporter.RepoUrl)
		_ = redisCli.HSet(key, "branch", buildExporter.Branch)
		_ = redisCli.HSet(key, "Sha1", buildExporter.Sha1)
	} else {
		// rebuild a test job: "build_cause":"MANUALTRIGGER,UPSTREAMTRIGGER", should !contains
		if !strings.Contains(build.BuildCause, "MANUALTRIGGER") {
			var upstreamJob, upstreamBuild string
			temp := strings.Split(buildExporter.BuildCause, `"`)
			if len(temp) != 3 {
				monitor.SendAlert(errors.New("unexpected build cause found: " + buildExporter.BuildCause))
			} else {
				if len(temp) > 1 {
					upstreamJob = temp[1]
				}
			}
			r, _ := regexp.Compile("([0-9]+)$")
			upstreamBuild = r.FindString(buildExporter.BuildCause)
			key = fmt.Sprintf("qex_%s_%s", upstreamJob, upstreamBuild)
			buildExporter.TestDetails.TestedRepo, _ = redisCli.HGet(key, "repo")
			buildExporter.TestDetails.TestedBranch, _ = redisCli.HGet(key, "branch")
			buildExporter.TestDetails.TestedSha1, _ = redisCli.HGet(key, "Sha1")
		}
	}

	changeItems := buildDetails.ChangeSet.Items
	r, _ := regexp.Compile("([A-Z]+)-([1-9][0-9]+)")
	for _, item := range changeItems {
		bugs := r.FindAllString(strings.ToUpper(item.Comment), 5)
		for _, bug := range bugs {
			buildExporter.Bugs[bug] = ""
		}
		bugs = r.FindAllString(strings.ToUpper(item.Msg), 5)
		for _, bug := range bugs {
			buildExporter.Bugs[bug] = ""
		}
	}

	bugs := r.FindAllString(strings.ToUpper(buildExporter.Branch), 5)
	for _, bug := range bugs {
		buildExporter.Bugs[bug] = ""
	}

	if !buildExporter.IsTestJob {
		bugs := make([]string, len(buildExporter.Bugs))

		i := 0
		for k, _ := range buildExporter.Bugs {
			bugs[i] = k
			i++
		}

		_ = redisCli.HSet(key, "bugs", strings.Join(bugs, ","))
	} else {
		str, _ := redisCli.HGet(key, "bugs")
		bugs := strings.Split(str, ",")
		for _, bug := range bugs {
			buildExporter.Bugs[bug] = ""
		}
	}

	UpdateDownstream(buildExporter)
}

func UpdateDownstream(buildExporter protocol.JenkinsBuildExporter) {
	delete(buildExporter.Bugs, "")
	log.Println("debug: buildExporter ###########")
	s, _ := json.MarshalIndent(buildExporter, "", "\t")
	fmt.Println(string(s))

	if producer, err := kafka.InitProducer(config.JenkinsBootstramp); err != nil {
		monitor.SendAlert(err)
	} else {
		bytes, _ := json.Marshal(buildExporter)
		kafka.Send(producer, config.JenkinsBuildExporterTopic, string(bytes))
	}

	client := jiraUtils.GetJiraClient()
	for bug, _ := range buildExporter.Bugs {
		jiraUtils.AddComment(client, bug, buildExporter)
	}

	influx_utils.ProcessJenkinsBuild(buildExporter)
}

func RetrieveBuildDetails(msg string) (build protocol.JenkinsBuild, buildDetails protocol.JenkinsBuildDetails) {
	json.Unmarshal([]byte(msg), &build)

	for {
		var resp *http.Response
		resp, err := http.Get(build.BuildUrl)
		if err != nil {
			monitor.SendAlert(err)
		}
		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			monitor.SendAlert(err)
		}
		json.Unmarshal(body, &buildDetails)
		if !buildDetails.Building {
			break
		} else {
			monitor.SendAlert(errors.New(fmt.Sprintf("%s:%s", buildDetails.URL, " seems stuck")))
			time.Sleep(60 * time.Second) // TODO: ideally should not have to handle this else... but let's see what will happen
		}
	}
	return build, buildDetails
}

func ExtractBugsFromMRTitle(jobName, jobBuild string) (bugs []string) {
	buildDetails := protocol.JenkinsBuildDetails{}

	url := fmt.Sprintf(config.JenkinsBuildInfoURL, jobName, jobBuild)
	log.Println("URL: " + url)
	var resp *http.Response
	resp, err := http.Get(url)
	if err != nil {
		monitor.SendAlert(err)
	}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		monitor.SendAlert(err)
	}
	json.Unmarshal(body, &buildDetails)

	var msg string
	for _, action := range buildDetails.Actions {
		if action.Class == "hudson.model.CauseAction" {
			if action.Causes[0].Class == "com.dabsquared.gitlabjenkins.cause.GitLabWebHookCause" {
				msg = action.Causes[0].ShortDescription
				break
			}
		}
	}
	fmt.Println(fmt.Sprintf("%s: %s: %s", jobName, jobBuild, msg))
	if msg == "" {
		return bugs
	}

	rMsg, _ := regexp.Compile("(/[0-9]+\")")
	iid := rMsg.FindString(msg) // /id"
	iid = strings.Replace(iid, "/", "", 1)
	iid = strings.Replace(iid, "\"", "", 1)
	mrIID, _ := strconv.Atoi(iid)

	var projectId int
	for k, v := range config.MapGitProductRepo {
		if strings.Contains(msg, k) {
			projectId = v
		}
	}
	if projectId == 0 {
		return bugs
	}

	git, _ := gitlab.NewClient(config.GitReadOnlyToken, gitlab.WithBaseURL(config.GitV4API))
	mr, _ := gitUtils.QueryGitlabMRByBranches(git, projectId, mrIID)

	log.Println("debug: MR title" + mr.Title)

	r, _ := regexp.Compile("([A-Z]+)-([1-9][0-9]+)")
	bugs = r.FindAllString(strings.ToUpper(mr.Title), 5)
	return bugs
}
