package hltool

import (
	"gopkg.in/gomail.v2"
)

// Message 内容
type Message struct {
	From 		string
	To 			[]string
	Cc  		[]string
	Subject 	string
	ContentType string
	Content 	string
	Attach 		string
}

// NewMessage 返回消息对象
// from: 发件人
// subject: 标题
// contentType: 内容的类型 text/plain text/html
// attach: 附件
// to: 收件人
// cc: 抄送人
func NewMessage(from, subject, contentType, content, attach string, to, cc []string) *Message{
	return &Message{
		From: from,
		Subject: subject,
		ContentType: contentType,
		Content: content,
		To: to,
		Cc: cc,
		Attach: attach,
	}
}

// Client 发送客户端
type Client struct {
	Host 	string
	Port	int
	Username string
	Password string
}

// NewClient 返回一个邮件客户端 
// host smtp地址
// username 用户名
// password 密码
// port 端口
func NewClient(host, username, password string, port int) *Client{
	return &Client{
		Host: host,
		Port: port,
		Username: username,
		Password: password,
	}
}

// Send 发送邮件
func (c *Client) Send(m *Message) (bool, error) {

	e := gomail.NewDialer(c.Host, c.Port, c.Username, c.Password)
	dm := gomail.NewMessage()
	dm.SetHeader("From", m.From)
	dm.SetHeader("To", m.To...)

	if len(m.Cc) != 0 {
		dm.SetHeader("Cc", m.Cc...)
	}

	dm.SetHeader("Subject", m.Subject)
	dm.SetBody(m.ContentType, m.Content)

	if m.Attach != "" {
		dm.Attach(m.Attach)
	}
	
	if err := e.DialAndSend(dm); err != nil {
		return false, err
	}
	return true, nil
}