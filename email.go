package hltool


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
}

// NewEmail 初始化 email
func NewEmail() *Email{
	return new(Email)
}

// Send 发送邮件
func (e *Email) Send(){
	email
}