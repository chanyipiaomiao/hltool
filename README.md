# hltool

Go 开发常用工具库


# 安装

使用golang官方 dep 管理依赖
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

- 支持按天分割日志
- 不同级别输出到不同文件
- 支持 文本/json日志类型,默认是json类型
- 设置日志最大保留时间

```go
import (
	"github.com/chanyipiaomiao/hltool"
)

func main() {
	
	commonFields := map[string]interface{}{
		"name": "zhangsan",
		"age":  "20",
	}

	hlog, _ := hltool.NewHLog("./test.log")
	// hlog.SetLevel("debug") debug|info|warn|error|fatal|panic
	logger, _ := hlog.GetLogger()

	// Info Warn 会输出到不同的文件
	logger.Info(commonFields, "测试Info消息")
	logger.Warn(commonFields, "测试Warn消息")
	
	// Error Fatal Panic 会输出到一个文件
	logger.Error(commonFields, "测试Error消息")
	logger.Fatal(commonFields, "测试Fatal消息")
	logger.Panic(commonFields, "测试Panic消息")
}
```
日志文件内容:
```shell
{"age":"20","level":"debug","msg":"测试Debug消息","name":"zhangsan","time":"2018-02-08 21:28:29"}
{"age":"20","level":"info","msg":"测试Info消息","name":"zhangsan","time":"2018-02-08 21:28:29"}
{"age":"20","level":"warning","msg":"测试Warn消息","name":"zhangsan","time":"2018-02-08 21:28:29"}
{"age":"20","level":"error","msg":"测试Error消息","name":"zhangsan","time":"2018-02-08 21:28:29"}
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