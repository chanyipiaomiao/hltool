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