package influx_utils

import (
	"github.com/luohuahuang/qex/config"
	"github.com/stretchr/testify/assert"
	"testing"
)

func Test_simpleWrite(t *testing.T) {
	err := WriteToInflux(config.TestExecutionMeasurement, map[string]string{
		"run_id": "hahahehe",
	}, map[string]interface{}{
		"product":          "chuky",
		"sub_product_line": "chuky",
		"service":          "chuky",
		"api":              "chuky",
		"case":             "chuky",
		"branch":           "chuky",
		"maintainer":       "chuky",
		"timestamp":        1,
		"duration":         1.0,
		"status":           1,
	})
	assert.NoError(t, err)
}
