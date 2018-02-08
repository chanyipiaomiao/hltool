# hltool

Go 开发常用工具库

# 依赖

```go
go get golang.org/x/crypto/scrypt
go get github.com/sirupsen/logrus
go get github.com/lestrrat/go-file-rotatelogs
go get github.com/rifflock/lfshook
go get gopkg.in/gomail.v2
```

# 安装
```go
go get github.com/chanyipiaomiao/hltool
```

# 钉钉机器人通知 示例
```go
import (
	"log"
	"github.com/chanyipiaomiao/hltool"
)

dingtalk := hltool.NewDingTalkClient("钉钉机器URL", "消息内容", "text|markdown")
ok, err := hltool.SendMessage(dingtalk)
if err != nil {
	log.Fatalf("发送钉钉通知失败了: %s", err)
}

```

# 发送邮件 示例
```go
import (
	"log"
	"github.com/chanyipiaomiao/hltool"
)

username := "xxxx@xxx.com"
host := "smtp.exmail.qq.com"
password := "password"
port := 465

subject := "主题"
content := "内容"
contentType := "text/plain|text/html"
attach := "附件路径" 或者 ""
to := []string{"xxx@xxx.com", "xxx@xx.com"}
cc := []string{"xxx@xxx.com", "xxx@xx.com"}

message := hltool.NewEmailMessage(username, subject, contentType, content, attach, to, cc)
email := hltool.NewEmailClient(host, username, password, port, message)
ok, err := hltool.SendMessage(email)
if err != nil {
	log.Fatalf("发送邮件失败了: %s", err)
}
```

# logrus Log库 示例
```go
import (
	"github.com/chanyipiaomiao/hltool"
)

func main() {
	hlog, _ := hltool.NewHLog("./", "test.log")
	logger, _ := hlog.GetLogger()
	contextLogger := logger.WithFields(hlog.GetLogField(map[string]interface{}{
		"username": "admin",
	}))
	contextLogger.Info("测试代码")
	contextLogger.Error("测试代码")
}
```
日志文件内容:
```shell
{"level":"info","msg":"测试代码","time":"2018-02-06 21:42:13","username":"admin"}
{"level":"error","msg":"测试代码","time":"2018-02-06 21:42:13","username":"admin"}
```

# BoltDB 嵌入式KV数据库 示例
```go
import (
	"log"

	"github.com/chanyipiaomiao/hltool"
)

func main() {

	// 数据库文件路径 表名
	db, err := hltool.NewBoltDB("./data/app.db", "token")
	if err != nil {
		log.Fatalf("%s", err)
	}
	db.Set(map[string][]byte{
		"hello": []byte("world"),
		"go":    []byte("golang"),
	})
	r, err := db.Get([]string{"hello", "go"})
	if err != nil {
		log.Fatalf("%s", err)
	}
	log.Println(r)
}
```