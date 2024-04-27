package ddding

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
)

func SendToDingTalkGroup(message string) {
	type DingTalkMessage struct {
		MsgType string `json:"msgtype"`
		Text    struct {
			Content string `json:"content"`
		} `json:"text"`
	}

	dingTalkURL := "https://oapi.dingtalk.com/robot/send?access_token=0cb614a04259e375d7a1e2a89879f06dbd8fca755d465598ca7e5673dcca2090"

	dingTalkMsg := DingTalkMessage{
		MsgType: "text",
	}
	dingTalkMsg.Text.Content = message

	jsonValue, _ := json.Marshal(dingTalkMsg)
	_, err := http.Post(dingTalkURL, "application/json", bytes.NewBuffer(jsonValue))
	if err != nil {
		fmt.Printf("消息发送失败:%v ", err)
	}
}
