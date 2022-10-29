package main

import (
	"fmt"
	"github.com/luohuahuang/qex/config"
	"github.com/luohuahuang/qex/internal/cache"
	gitUtils "github.com/luohuahuang/qex/internal/git"
	influxUtils "github.com/luohuahuang/qex/internal/influx"
	"github.com/luohuahuang/qex/pkg/mattermost"
	"github.com/luohuahuang/qex/protocol"
	"log"
	"time"

	"github.com/xanzy/go-gitlab"
)

func main() {
	git, err := gitlab.NewClient(config.GitReadOnlyToken, gitlab.WithBaseURL(config.GitV4API))
	if err != nil {
		mattermost.SendAlert(err, config.MatterMostMonitor)
	}

	now := time.Now()
	year := now.Year() - 1
	// our fiscal year start from 1st Oct.
	if now.Month() == time.October || now.Month() == time.November || now.Month() == time.December {
		year = now.Year()
	}

	startTime := time.Date(year, time.October, 1, 0, 0, 0, 0, time.Local)

	log.Println(now.Year())

	for k, v := range config.MapGitTestRepo {
		runId := fmt.Sprintf("%s", time.Now().Format("2006-01-02-15:04:05"))

		mrs, _ := gitUtils.QueryGitlabProjectMRs(git, v, &startTime, &now)
		for _, mr := range mrs {
			gitMR := protocol.GitMR{
				RunId:   runId,
				Product: k,
				MrID:    mr.IID,
				Author:  mr.Author.Username,
				State:   mr.State,
			}
			influxUtils.ProcessGitMR(gitMR)
		}
		log.Println(fmt.Sprintf("%s: %d", k, len(mrs)))
	}
}

// TODO: https://stackoverflow.com/questions/35373995/github-user-email-is-null-despite-useremail-scope
func QueryEmailByUserId(client *gitlab.Client, id int) string {
	redisCli := cache.New(config.CacheServer)
	email, err := redisCli.Get(fmt.Sprintf(config.GitUserCacheFormat, id))
	if err != nil {
		mattermost.SendAlert(err, config.MatterMostMonitor)
	}
	if email != "" {
		return email
	}
	opt := &gitlab.GetUsersOptions{}
	user, resp, err := client.Users.GetUser(id, *opt)
	if err != nil {
		log.Panic(err.Error())
	}
	log.Println(resp.Status)
	redisCli.Set(fmt.Sprintf(config.GitUserCacheFormat, id), user.Email, 0)
	return user.Email
}
