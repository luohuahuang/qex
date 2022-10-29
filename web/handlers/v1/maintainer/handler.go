package maintainer

import (
	"encoding/json"
	"github.com/gin-gonic/gin"
	"github.com/luohuahuang/qex/config"
	"github.com/luohuahuang/qex/internal/maintainer"
	"github.com/luohuahuang/qex/pkg/mattermost"
	"github.com/luohuahuang/qex/protocol"
	"net/http"
)

func Handler(c *gin.Context) {
	var cases protocol.Cases

	if err := c.BindJSON(&cases); err != nil {
		mattermost.SendAlert(err, config.MatterMostMonitor)
		return
	}

	info := maintainer.Process(cases)
	JsonInfo, err := json.Marshal(info)
	if err != nil {
		mattermost.SendAlert(err, config.MatterMostMonitor)
		return
	}
	c.JSON(http.StatusOK, gin.H{"code": http.StatusOK,
		"message": string(JsonInfo)})
}
