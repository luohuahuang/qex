package mattermost

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/luohuahuang/qex/config"
	"github.com/luohuahuang/qex/protocol"
	"log"
	"net/http"
	"regexp"
	"strings"
)

type Msg struct {
	Username string `json:"username"`
	Text     string `json:"text"`
}

func SendAlert(err error, hook string) {
	log.Println(fmt.Sprintf("send error msg to mattermost: %s", err.Error()))
	msg := Msg{
		Text: fmt.Sprintf(":fire: :fire: :fire: %s", err.Error()),
	}
	send(msg, hook)
}

func send(msg Msg, hook string) {
	if msg.Username == "" {
		msg.Username = "QEX Monitor"
	}
	payload, _ := json.Marshal(msg)
	if _, err := http.Post(hook, "application/json", bytes.NewBuffer(payload)); err != nil {
		log.Println(err.Error())
	}
}

func SendMsg(buildExporter protocol.JenkinsBuildExporter, issueId ...string) {
	users := strings.Split(buildExporter.User, "@")
	if len(users) < 2 { // means not email address
		users = []string{"CI pipeline or System"}
	}

	var emojiResult string
	switch buildExporter.Result {
	case "SUCCESS":
		emojiResult = "✅"
	case "UNSTABLE":
		emojiResult = "⚠️"
	case "ABORTED":
		emojiResult = "⏹"
	default:
		emojiResult = "❌"
	}

	var msg Msg

	if buildExporter.IsTestJob {
		if len(issueId) > 0 { // send for qex-msg-bot
			msg = Msg{
				Text: fmt.Sprintf("INFO: QEX just commented [%s](https://jira.example.com/browse/%s). Test result: %s, Pass: %d, Fail: %d, Skipped: %d. Triggered by: @%s",
					issueId[0], issueId[0],
					emojiResult, buildExporter.TestDetails.TotalCount-buildExporter.TestDetails.FailCount-buildExporter.TestDetails.SkipCount, buildExporter.TestDetails.FailCount, buildExporter.TestDetails.SkipCount,
					users[0]),
			}
			send(msg, config.MatterMostPublic)
			return
		}
		msg = Msg{
			Text: fmt.Sprintf("Jenkins QEX: [%s](%s). Test result: %s, Pass: %d, Fail: %d, Skipped: %d. Triggered by: @%s",
				buildExporter.JobName,
				buildExporter.BuildUrl,
				emojiResult, buildExporter.TestDetails.TotalCount-buildExporter.TestDetails.FailCount-buildExporter.TestDetails.SkipCount, buildExporter.TestDetails.FailCount, buildExporter.TestDetails.SkipCount,
				users[0]),
		}
		log.Println("send mattermost msg for jenkins test job")
		for k, v :=range config.MatterMostEndpoints {
			r, _ := regexp.Compile(k)
			if r.FindString(buildExporter.JobName) != "" {
				send(msg, v)
			}
		}
	} else {
		if len(issueId) > 0 { // send for qex-msg-bot
			msg = Msg{
				Text: fmt.Sprintf("INFO: QEX just commented [%s](https://jira.example.com/browse/%s). Build result: %s. Triggered by: @%s",
					issueId[0], issueId[0],
					emojiResult,
					users[0]),
			}
			send(msg, config.MatterMostPublic)
			return
		}
		msg = Msg{
			Text: fmt.Sprintf("Jenkins QEX: [%s](%s). Build result: %s. Triggered by: @%s",
				buildExporter.JobName,
				buildExporter.BuildUrl,
				emojiResult,
				users[0]),
		}
		log.Println("send mattermost msg for jenkins build job")
		for k, v :=range config.MatterMostEndpoints {
			r, _ := regexp.Compile(k)
			if r.FindString(buildExporter.JobName) != "" {
				send(msg, v)
			}
		}
	}
}