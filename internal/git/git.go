package git_utils

import (
	"github.com/luohuahuang/qex/config"
	influxUtils "github.com/luohuahuang/qex/internal/influx"
	"github.com/luohuahuang/qex/monitor"
	"github.com/luohuahuang/qex/protocol"
)

func Process(git protocol.Git) {
	if err := influxUtils.WriteToInflux(config.GitMeasurement, map[string]string{
		"run_id": git.RunId,
	}, map[string]interface{}{
		"case":       git.Case,
		"maintainer": git.Maintainer,
		"product":    git.Product,
	}); err != nil {
		monitor.SendAlert(err)
	}
}
