package hltool

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net/http"
)

// DingTalkMessage 消息
type DingTalkMessage struct {
	Message string //消息
	Title   string // markdown标题
	Type    string // 消息类型
}

// DingTalkClient 通过钉钉机器人发送消息
type DingTalkClient struct {
	RobotURL string
	Message  *DingTalkMessage
}

// SendMessage 通过钉钉机器人发送消息
func (d *DingTalkClient) SendMessage() (bool, error) {

	var message string
	switch d.Message.Type {
	case "text":
		message = fmt.Sprintf(`{"msgtype": "text","text": {"content": "%s"}}`, d.Message.Message)
	case "markdown":
		message = fmt.Sprintf(`{"msgtype": "markdown","markdown":{"title": "%s", "text": "%s"}}`, d.Message.Title, d.Message.Message)
	default:
		message = fmt.Sprintf(`{"msgtype": "text","text": {"content": "%s"}}`, d.Message.Message)
	}

	client := &http.Client{}
	request, _ := http.NewRequest("POST", d.RobotURL, bytes.NewBuffer([]byte(message)))
	request.Header.Set("Content-type", "application/json")
	response, err := client.Do(request)
	if err != nil {
		return false, fmt.Errorf("访问钉钉URL(%s) 出错了: %s", d.RobotURL, err)
	}
	if response.StatusCode != 200 {
		body, _ := ioutil.ReadAll(response.Body)
		return false, fmt.Errorf("访问钉钉URL(%s) 出错了: %s", d.RobotURL, string(body))
	}
	ioutil.ReadAll(response.Body)
	return true, nil
}
