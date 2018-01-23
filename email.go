package hltool

import (
	"github.com/jordan-wright/email"
	"net/smtp"
)


// Auth 验证
type Auth struct {
	Host 	string
	Port	string
	Username string
	Password string
	Identity string
}

// Email 发送邮件
type Email struct {
	From 	string
	To 		[]string
	Bcc 	[]string
	Cc  	[]string
	Subject string
	Text 	[]byte
	HTML 	[]byte
	Auth    *Auth
}

// NewEmail 初始化 email
func NewEmail() *Email{
	return new(Email)
}

// Send 发送邮件
func (e *Email) Send(){
	m := email.NewEmail()
	m.From = e.From
	m.To = e.To
	m.Subject = e.Subject

	if len(e.Text) != 0 {
		m.Text = e.Text
	}

	if len(e.HTML) != 0 {
		m.HTML = e.HTML
	}

	m.Send(e.Auth.Host + ":" + e.Auth.Port, smtp.PlainAuth(e.Auth.Identity, e.Auth.Username, e.Auth.Password, e.Auth.Host))

	return
}