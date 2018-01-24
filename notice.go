package hltool

// Notice 通知接口
type Notice interface {
	SendMessage() (bool, error)
}

// SendMessage 发送消息
func SendMessage(notice Notice)(bool, error){
	return notice.SendMessage()
}
