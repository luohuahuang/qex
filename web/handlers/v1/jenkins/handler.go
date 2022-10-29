package jenkins

import (
	"encoding/json"
	"github.com/gin-gonic/gin"
	"github.com/luohuahuang/qex/config"
	"github.com/luohuahuang/qex/internal/kafka"
	"github.com/luohuahuang/qex/pkg/mattermost"
	"github.com/luohuahuang/qex/protocol"
	"log"
)

func Handler(c *gin.Context) {
	var build protocol.JenkinsBuild

	if err := c.BindJSON(&build); err != nil {
		mattermost.SendAlert(err, config.MatterMostMonitor)
		return
	} else {
		if producer, err := kafka.InitProducer(config.JenkinsBootstramp); err != nil {
			mattermost.SendAlert(err, config.MatterMostMonitor)
		} else {
			bytes, _ := json.Marshal(build)
			log.Println(string(bytes))
			kafka.Send(producer, config.JenkinsBuildTopic, string(bytes))
		}
	}
}
