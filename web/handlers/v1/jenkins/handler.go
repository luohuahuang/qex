package jenkins

import (
	"encoding/json"
	"github.com/luohuahuang/qex/config"
	"github.com/luohuahuang/qex/internal/kafka"
	"github.com/luohuahuang/qex/monitor"
	"github.com/luohuahuang/qex/protocol"
	"github.com/gin-gonic/gin"
	"log"
)

func Handler(c *gin.Context) {
	var build protocol.JenkinsBuild

	if err := c.BindJSON(&build); err != nil {
		monitor.SendAlert(err)
		return
	} else {
		if producer, err := kafka.InitProducer(config.JenkinsBootstramp); err != nil {
			monitor.SendAlert(err)
		} else {
			bytes, _ := json.Marshal(build)
			log.Println(string(bytes))
			kafka.Send(producer, config.JenkinsBuildTopic, string(bytes))
		}
	}
}
