package influx_utils

import (
	"context"
	influx "github.com/influxdata/influxdb-client-go/v2"
	"github.com/luohuahuang/qex/config"
	"github.com/luohuahuang/qex/monitor"
	"time"
)

func WriteToInflux(measurement string, tags map[string]string, fields map[string]interface{}) error {
	ctx := context.Background()
	client := influx.NewClient(config.InfluxUrl, config.InfluxToken)
	defer client.Close()
	writeClient := client.WriteAPIBlocking(config.InfluxOrg, config.InfluxBucket)
	p := influx.NewPoint(measurement,
		tags,
		fields,
		time.Now())
	err := writeClient.WritePoint(ctx, p)
	if err != nil {
		monitor.SendAlert(err)
		return err
	}
	return nil
}
