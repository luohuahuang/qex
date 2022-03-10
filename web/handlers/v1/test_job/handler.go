package test_job

import (
	"github.com/luohuahuang/qex/internal/test_execution"
	"github.com/luohuahuang/qex/monitor"
	"github.com/luohuahuang/qex/protocol"
	"github.com/gin-gonic/gin"
)

func Handler(c *gin.Context) {
	var testCase protocol.QEXTestCase

	if err := c.BindJSON(&testCase); err != nil {
		monitor.SendAlert(err)
		return
	}

	if err := test_execution.Process(testCase); err != nil {
		monitor.SendAlert(err)
		return
	}
}
