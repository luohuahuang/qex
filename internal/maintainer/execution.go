package maintainer

import (
	"github.com/luohuahuang/qex/config"
	"github.com/luohuahuang/qex/internal/cache"
	"github.com/luohuahuang/qex/pkg/mattermost"
	"github.com/luohuahuang/qex/protocol"
)

func Process(cases protocol.Cases) (info protocol.Maintainers) {
	info = protocol.Maintainers{
		Data: map[string]string{},
	}
	redisCli := cache.New(config.CacheServer)
	for _, c := range cases.Data {
		value, err := redisCli.Get(c)
		if err != nil {
			mattermost.SendAlert(err, config.MatterMostMonitor)
		}
		info.Data[c] = value
	}
	return info
}
