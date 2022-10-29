package xml_harvestor

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/luohuahuang/qex/config"
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

const (
	XmlFileKey   = "XML_RESULTS"
	ClassNameKey = "PACKAGE_NAME"
)

func Handler(c *gin.Context) {
	t := time.Now()
	msgId := fmt.Sprintf("%d-%02d-%02dT%02d:%02d:%02d",
		t.Year(), t.Month(), t.Day(),
		t.Hour(), t.Minute(), t.Second(),
	)

	kv := extraFormParams(c)
	pkgName, found := kv[ClassNameKey]
	if !found {
		err := fmt.Errorf("did not find param %s", ClassNameKey)
		mattermost.SendAlert(err, config.MatterMostMonitor)
		c.JSON(http.StatusInternalServerError, gin.H{
			"fail":  err.Error(),
			"msgId": msgId,
		})
		return
	}
	log.Printf("processing request for job: %s (%s)", pkgName, msgId)

	file, err := c.FormFile(XmlFileKey)
	if err != nil {
		err = errors.Wrap(err, fmt.Sprintf("error parsing %s", XmlFileKey))
		mattermost.SendAlert(err, config.MatterMostMonitor)
		c.JSON(http.StatusInternalServerError, gin.H{
			"fail":  err.Error(),
			"msgId": msgId,
		})
		return
	}

	dst, err := upload(file, pkgName)
	if err != nil {
		err = errors.Wrap(err, "error uploading file")
		mattermost.SendAlert(err, config.MatterMostMonitor)
		c.JSON(http.StatusInternalServerError, gin.H{
			"fail":  err.Error(),
			"msgId": msgId,
			"dst":   dst,
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"success": "uploaded successfully",
		"msgId":   msgId,
		"dst":     dst,
	})
}

func extraFormParams(c *gin.Context) map[string]string {
	envs := strings.Split(c.PostForm("params"), "\n")
	kv := map[string]string{}
	for _, env := range envs {
		arr := strings.Split(env, "=")
		if len(arr) > 1 {
			kv[arr[0]] = arr[1]
		}
	}
	return kv
}

func upload(file *multipart.FileHeader, pkgName string) (dst string, err error) {
	src, err := file.Open()
	if err != nil {
		return "", errors.Wrapf(err, fmt.Sprintf("error opening file: %s", src))
	}
	defer src.Close()

	product := strings.SplitN(pkgName, ".", 2)[0]
	folder := fmt.Sprintf("%s/all-report-for-%s", config.XmlFileDir, product)
	if _, err := os.Stat(folder); os.IsNotExist(err) {
		err = os.MkdirAll(folder, 0755)
		if err != nil {
			return "", errors.Wrapf(err, fmt.Sprintf("unable to make path: %s", folder))
		}
		log.Printf("Created directory %s", folder)
		mattermost.SendAlert(fmt.Errorf("created folder %s", folder), config.MatterMostMonitor)
	}

	filename := strings.SplitN(file.Filename, ".xml", 2)[0]

	dst = fmt.Sprintf("%s/%s.xml", folder, filename)
	out, err := os.Create(dst)
	if err != nil {
		return dst, errors.Wrapf(err, fmt.Sprintf("error creating dst: %s", dst))
	}
	defer out.Close()

	_, err = io.Copy(out, src)
	if err != nil {
		return dst, errors.Wrapf(err, "error copying file")
	}
	return
}
