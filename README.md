# hltool

常用工具箱

# 依赖

```go
go get golang.org/x/crypto/scrypt
go get github.com/sirupsen/logrus
go get github.com/lestrrat/go-file-rotatelogs
go get github.com/rifflock/lfshook
go get gopkg.in/gomail.v2
```

# Log 示例
```go

import (
	"log"

	"github.com/chanyipiaomiao/hltool"
)

func main() {
	hlog, err := hltool.NewHLog("./log", "test.log")
	if err != nil {
		log.Fatalf("%s", err)
	}
	hlog.SetCommonFields(map[string]interface{}{
		"haha":  "这是公共字段",
		"stime": hltool.GetNowTime2(),
	})
	logger, err := hlog.GetLogger()
	if err != nil {
		log.Fatalf("%s", err)
	}
	logger.Info("测试一下")
	logger.Warn("测试一下")
	logger.Error("测试一下")

}

```