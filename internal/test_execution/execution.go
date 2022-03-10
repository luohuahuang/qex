package test_execution

import (
	"github.com/luohuahuang/qex/config"
	influxUtils "github.com/luohuahuang/qex/internal/influx"
	"github.com/luohuahuang/qex/protocol"
)

func Process(testCase protocol.QEXTestCase) error {
	if err := influxUtils.WriteToInflux(config.TestExecutionMeasurement, map[string]string{
		"run_id": testCase.RunId,
	}, map[string]interface{}{
		"product":          testCase.Product,
		"sub_product_line": testCase.SubProduct,
		"service":          testCase.Service,
		"api":              testCase.API,
		"case":             testCase.Case,
		"branch":           testCase.Branch,
		"maintainer":       testCase.Maintainer,
		"timestamp":        testCase.Timestamp,
		"duration":         testCase.Duration,
		"status":           testCase.Status,
	}); err != nil {
		return err
	}
	return nil
}
