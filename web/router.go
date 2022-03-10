package web

import (
	"github.com/luohuahuang/qex/web/handlers/v1/git"
	"github.com/luohuahuang/qex/web/handlers/v1/test_job"
	"github.com/gin-contrib/pprof"
	"github.com/gin-gonic/gin"
	"os"
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
	testRunAPIs.POST("/case", test_job.Handler)

	gitAPIs := v1API.Group("/git")
	gitAPIs.POST("/upload", git.Handler)

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
