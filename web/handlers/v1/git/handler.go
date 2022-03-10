package git

import (
	"fmt"
	"github.com/luohuahuang/qex/config"
	"github.com/luohuahuang/qex/internal/kafka"
	"github.com/luohuahuang/qex/monitor"
	"github.com/gin-gonic/gin"
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
	log.Println("processing request for job: " + kv["PRODUCT"])

	t := time.Now()
	msgId := fmt.Sprintf("%d-%02d-%02dT%02d:%02d:%02d-%s",
		t.Year(), t.Month(), t.Day(),
		t.Hour(), t.Minute(), t.Second(), kv["PRODUCT"],
	)

	if file, err := c.FormFile(config.GitMasterReport); err != nil {
		monitor.SendAlert(err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"fail": err.Error(),
		})
		return
	} else {
		if err := upload(file, msgId); err != nil {
			monitor.SendAlert(err)
			return
		} else {
			if producer, err := kafka.InitProducer(config.GitMasterBootstramp); err != nil {
				monitor.SendAlert(err)
			} else {
				kafka.Send(producer, config.GitMasterTopic, msgId)
			}

		}
	}
	c.JSON(http.StatusOK, gin.H{
		"success": "uploaded successfully",
	})
}

func extraFormParams(c *gin.Context) map[string]string {
	envs := strings.Split(c.PostForm("params"), "\n")
	kv := map[string]string{}
	for _, env := range envs {
		arr := strings.Split(env, "=")
		if len(arr) > 1 {
			//log.Println("#" + arr[0] + "#" + arr[1] + "#")
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

	return err
}
