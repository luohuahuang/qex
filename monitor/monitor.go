package monitor

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/luohuahuang/qex/config"
	"log"
	"net/http"
)

type Msg struct {
	Username string `json:"username"`
	Text     string `json:"text"`
}

func SendAlert(err error) {
	log.Println(fmt.Sprintf("send error msg to mattermost: %s", err.Error()))
	msg := Msg{
		Username: "QEX Monitor",
		Text:     fmt.Sprintf(":fire: :fire: :fire: %s", err.Error()),
	}
	payload, _ := json.Marshal(msg)
	if _, err := http.Post(config.MsgServer, "application/json", bytes.NewBuffer(payload)); err != nil {
		log.Println(err.Error())
	}
}
