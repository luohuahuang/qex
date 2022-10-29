package test_job

import (
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/luohuahuang/qex/config"
	"github.com/luohuahuang/qex/internal/test_execution"
	"github.com/luohuahuang/qex/pkg/mattermost"
	"github.com/luohuahuang/qex/protocol"
)

func Handler(c *gin.Context) {
	var testCase protocol.QEXTestCase

	if err := c.BindJSON(&testCase); err != nil {
		mattermost.SendAlert(err, config.MatterMostMonitor)
		return
	}

	if err := test_execution.Process(testCase); err != nil {
		mattermost.SendAlert(errors.New(fmt.Sprintf("[ERROR] Test Execution: %v\n%v", err.Error(), testCase)), config.MatterMostMonitor)
		return
	}
}
