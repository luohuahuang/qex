package main

import (
	"bufio"
	"errors"
	"fmt"
	"github.com/luohuahuang/qex/config"
	"github.com/luohuahuang/qex/internal/influx"
	"github.com/luohuahuang/qex/internal/kafka"
	"github.com/luohuahuang/qex/monitor"
	"github.com/luohuahuang/qex/protocol"
	"log"
	"os"
	"strings"
	"time"
)

func ProcessEvent() {
	Consume(kafka.New(config.GitMasterTopic, config.GitMasterZooKeeper))
}

func main() {
	ProcessEvent()
}

func Consume(k *kafka.Consumer) {
	for {
		select {
		case msg := <-k.ConsumerGroup.Messages():
			if msg == nil || msg.Topic != k.Topic {
				continue
			}
			log.Println(fmt.Sprintf("topic: %s; msg: %s", msg.Topic, string(msg.Value)))

			process(string(msg.Value))
			err := k.ConsumerGroup.CommitUpto(msg)
			if err != nil {
				log.Println("error commit zookeeper: ", err.Error())
			}
		}
	}
}

func process(msgId string) {
	arr := strings.Split(msgId, "-")
	product := arr[len(arr)-1]
	if file, err := os.Open(fmt.Sprintf(config.GitMasterLogFullPathFormat, msgId)); err != nil {
		monitor.SendAlert(err)
	} else {
		scanner := bufio.NewScanner(file)
		for scanner.Scan() {
			// ea70897 keshia.yapkt@example.com 2022-02-16 Test_abc_def_verify_ghi 18 ./get_param_test.go
			row := strings.Split(scanner.Text(), " ")
			if len(row) < 4 {
				monitor.SendAlert(errors.New(fmt.Sprintf("found corrupted git record: %s", row)))
			} else {
				git := protocol.Git{
					RunId:      fmt.Sprintf("%s", time.Now().Format("2006-01-02-15:04:05")),
					CommitId:   row[0],
					Maintainer: row[1],
					Case:       row[3],
					Product:    product,
				}
				influx_utils.ProcessGitMaintainer(git)
			}
		}
	}

}
