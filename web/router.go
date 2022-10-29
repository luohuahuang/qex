package web

import (
	"github.com/gin-contrib/pprof"
	"github.com/gin-contrib/timeout"
	"github.com/gin-gonic/gin"
	"github.com/luohuahuang/qex/web/handlers/v1/git"
	"github.com/luohuahuang/qex/web/handlers/v1/jenkins"
	"github.com/luohuahuang/qex/web/handlers/v1/maintainer"
	"github.com/luohuahuang/qex/web/handlers/v1/test_job"
	"github.com/luohuahuang/qex/web/handlers/v1/xml_harvestor"
	"os"
	"time"
)

type Router struct {
	*gin.Engine
}

func NewRouter() *Router {
	gin.SetMode(gin.DebugMode)

	router := gin.New()

	pprof.Register(router)
	coreAPI := router.Group("/qex")
	v1API := coreAPI.Group("/v1")

	testRunAPIs := v1API.Group("/test_job")
	testRunAPIs.POST("/case", timeout.New(
		timeout.WithTimeout(100*time.Second),
		timeout.WithHandler(test_job.Handler),
	))

	maintainerAPIs := v1API.Group("/maintainer")
	maintainerAPIs.GET("/", timeout.New(
		timeout.WithTimeout(100*time.Second),
		timeout.WithHandler(maintainer.Handler),
	))

	gitAPIs := v1API.Group("/git")
	gitAPIs.POST("/upload", timeout.New(
		timeout.WithTimeout(100*time.Second),
		timeout.WithHandler(git.Handler),
	))

	harvestorAPIs := v1API.Group("/harvestor")
	harvestorAPIs.POST("/upload", timeout.New(
		timeout.WithTimeout(100*time.Second),
		timeout.WithHandler(xml_harvestor.Handler),
	))

	jenkinsAPIs := v1API.Group("/jenkins")
	jenkinsAPIs.POST("/upload", timeout.New(
		timeout.WithTimeout(100*time.Second),
		timeout.WithHandler(jenkins.Handler),
	))

	return &Router{
		router,
	}
}

func Run() {
	if port := os.Getenv("PORT"); port != "" {
		NewRouter().Run(port)
	} else {
		NewRouter().Run(":8888")
	}
}
