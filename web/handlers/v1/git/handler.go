package git

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/luohuahuang/qex/config"
	"github.com/luohuahuang/qex/internal/kafka"
	"github.com/luohuahuang/qex/pkg/mattermost"
	"github.com/pkg/errors"
	"io"
	"log"
	"mime/multipart"
	"net/http"
	"os"
	"strings"
	"time"
)

func Handler(c *gin.Context) {
	kv := extraFormParams(c)
	product, found := kv["PRODUCT"]
	if !found {
		err := errors.New("PRODUCT param not found")
		mattermost.SendAlert(err, config.MatterMostMonitor)
		c.JSON(http.StatusInternalServerError, gin.H{
			"fail": err.Error(),
		})
		return
	}
	log.Println("processing request for job: " + product)

	t := time.Now()
	msgId := fmt.Sprintf("%d-%02d-%02dT%02d:%02d:%02d-%s",
		t.Year(), t.Month(), t.Day(),
		t.Hour(), t.Minute(), t.Second(), product,
	)

	file, err := c.FormFile(config.GitMasterReport)
	if err != nil {
		err = errors.Wrap(err, "error adding form file")
		mattermost.SendAlert(err, config.MatterMostMonitor)
		c.JSON(http.StatusInternalServerError, gin.H{
			"fail": err.Error(),
		})
		return
	}
	if err := upload(file, msgId); err != nil {
		err = errors.Wrap(err, "error uploading file")
		mattermost.SendAlert(err, config.MatterMostMonitor)
		c.JSON(http.StatusInternalServerError, gin.H{
			"fail": err.Error(),
		})
		return
	}

	producer, err := kafka.InitProducer(config.GitMasterBootstramp)
	if err != nil {
		mattermost.SendAlert(errors.Wrap(err, "error in kafka.InitProducer"), config.MatterMostMonitor)
		c.JSON(http.StatusInternalServerError, gin.H{
			"fail": err.Error(),
		})
		return
	}

	if err := kafka.Send(producer, config.GitMasterTopic, msgId); err != nil {
		mattermost.SendAlert(errors.Wrap(err, "error in kafka.Send"), config.MatterMostMonitor)
		c.JSON(http.StatusInternalServerError, gin.H{
			"fail": err.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"success": "uploaded successfully",
	})
}

func extraFormParams(c *gin.Context) map[string]string {
	envs := strings.Split(c.PostForm("params"), "\n")
	kv := map[string]string{}
	for _, env := range envs {
		arr := strings.SplitN(env, "=", 2)
		if len(arr) > 1 {
			kv[arr[0]] = arr[1]
		}
	}
	return kv
}

func upload(file *multipart.FileHeader, msgId string) (err error) {
	src, err := file.Open()
	if err != nil {
		return err
	}
	defer src.Close()

	dst := fmt.Sprintf(config.GitMasterLogFullPathFormat, msgId)
	out, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer out.Close()

	_, err = io.Copy(out, src)
	if err != nil {
		return
	}
	return
}
