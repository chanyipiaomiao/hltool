# hltool

Go 开发常用工具库


# 安装

使用golang官方 dep 管理依赖
```go
go get github.com/chanyipiaomiao/hltool
```

# 功能列表
- [AES加密解密](#aes加密解密)
- [RSA加密解密](#rsa加密解密)
- [钉钉机器人通知](#钉钉机器人通知)
- [发送邮件](#发送邮件)
- [JWT Token生成解析](#jwt-token生成解析)
- [Log库](#log库)
- [BoltDB嵌入式KV数据库](#boltdb嵌入式kv数据库)
- [检测图片类型](#检测图片类型)
- [图片转[]byte](#图片转byte数组)
- [[]byte转换为png/jpg](#byte数组转换为png-jpg)
- [json文件转换为byte数组](#json文件转换为byte数组)
- [json []byte转换为struct](#json-byte数组转换为-struct)
- [struct序列化成二进制文件和反序列化](#struct序列化成二进制文件和反序列化)
- [struct序列化成byte数组和反序列化](#struct序列化成byte数组和反序列化)

### AES加密解密

```go
package main

import (
	"encoding/base64"
	"fmt"
	"log"

	"github.com/chanyipiaomiao/hltool"
)

func main() {

	// AES 加解密 指定加密的密码
	goaes := hltool.NewGoAES([]byte("O8Hp8WQbFPT7b5AUsEMVLtIU3MVYOrt8"))

	// 加密数据
	encrypt, err := goaes.Encrypt([]byte("123456"))
	if err != nil {
		log.Fatalln(err)
	}

	fmt.Println(base64.StdEncoding.EncodeToString(encrypt))

	// 解密数据
	decrypt, err := goaes.Decrypt(encrypt)
	if err != nil {
		log.Fatalln(err)
	}
	fmt.Println(string(decrypt))

}

```

[返回到目录](#功能列表)

### RSA加密解密

```go
package main

import (
	"fmt"
	"log"
	"github.com/chanyipiaomiao/hltool"
)
func main() {

	// 生成 2048 位密钥对文件 指定名称
	err := hltool.NewRSAFile("id_rsa.pub", "id_rsa", 2048)
	if err != nil {
		log.Fatalln(err)
	}

	// 生成密钥对字符串
	// pub, pri, err := hltool.NewRSAString(2048)
	// if err != nil {
	// 	log.Fatalln(err)
	// }
	// fmt.Println(pub)
	// fmt.Println(pri)

	// 指定 公钥文件名 和 私钥文件名
	gorsa, err := hltool.NewGoRSA("id_rsa.pub", "id_rsa")
	if err != nil {
		log.Fatalln(err)
	}

	// 明文字符
	rawStr := "O8Hp8WQbFPT7b5AUsEMVLtIU3MVYOrt8"

	// 使用公钥加密
	encrypt, err := gorsa.PublicEncrypt([]byte(rawStr))
	if err != nil {
		log.Fatalln(err)
	}

	// 使用私钥解密
	decrypt, err := gorsa.PrivateDecrypt(encrypt)
	if err != nil {
		log.Fatalln(err)
	}
	fmt.Println(string(decrypt))
}
```
[返回到目录](#功能列表)

### 钉钉机器人通知
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
[返回到目录](#功能列表)

### 发送邮件
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
[返回到目录](#功能列表)

### JWT Token生成解析
```go
import (
	"fmt"
	"log"

	"github.com/chanyipiaomiao/hltool"
)

func main() {

	// 签名字符串
	sign := "fDEtrkpbQbocVxYRLZrnkrXDWJzRZMfO"

	token := hltool.NewJWToken(sign)

	// -----------  生成jwt token -----------
	tokenString, err := token.GenJWToken(map[string]interface{}{
		"name": "root",
	})
	if err != nil {
		log.Fatalf("%s", err)
	}
	fmt.Println(tokenString)

	// -----------  解析 jwt token -----------
	r, err := token.ParseJWToken(tokenString)
	if err != nil {
		log.Fatalf("%s", err)
	}
	fmt.Println(r)

}

```
输出
```shell
eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJuYW1lIjoicm9vdCJ9.NJMXxkzdBBWrNUO5u2oXFLU9FD18TWiXHqxM2msT1x0

map[name:root]
```
[返回到目录](#功能列表)

### Log库

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
[返回到目录](#功能列表)

### BoltDB嵌入式KV数据库
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
[返回到目录](#功能列表)


### 检测图片类型

```go
package main

import (
	"fmt"

	"github.com/chanyipiaomiao/hltool"
)

func main() {

	bytes, _ := hltool.ImageToBytes("1.png")
	fmt.Println(hltool.ImageType(bytes))

}
```
输出结果:

```go
image/png
```

[返回到目录](#功能列表)

### 图片转byte数组

```go
package main

import (
	"fmt"

	"github.com/chanyipiaomiao/hltool"
)

func main() {

	bytes, err := hltool.ImageToBytes("1.png")
	if err != nil {
		fmt.Println(err)
	}

}
```

[返回到目录](#功能列表)

### byte数组转换为png jpg
```go
package main

import (
	"fmt"

	"github.com/chanyipiaomiao/hltool"
)

func main() {

	bytes, err := hltool.ImageToBytes("1.png")
	if err != nil {
		log.Fatalln(err)
	}

	err = hltool.BytesToImage(bytes, "111.png")
	if err != nil {
		log.Fatalln(err)
	}

}
```

[返回到目录](#功能列表)

### json文件转换为byte数组

json文件内容
```sh
{
    "Name": "张三",
    "Age": 20,
    "Address": {
        "Country": "China",
        "Province": "Shanghai",
        "City": "Shanghai"
}
```

```go
package main

import (
	"fmt"
	"log"

	"github.com/chanyipiaomiao/hltool"
)

func main() {

	// 读取json文件转换为 []byte
	b, err := hltool.JSONFileToBytes("/Users/helei/Desktop/test.json")
	if err != nil {
		log.Fatalln(err)
	}
}

```

[返回到目录](#功能列表)

### json byte数组转换为 struct

```go
package main

import (
	"fmt"
	"log"

	"github.com/chanyipiaomiao/hltool"
)

type Person struct {
	Name    string `json:"Name"`
	Age     int    `json:"Age"`
	Address struct {
		Country  string `json:"Country"`
		Province string `json:"Province"`
		City     string `json:"City"`
	} `json:"Address"`
}

func main() {

	// 读取json文件转换为 []byte
	b, err := hltool.JSONFileToBytes("/Users/helei/Desktop/test.json")
	if err != nil {
		log.Fatalln(err)
	}

	// json []byte转换为 struct
	p := new(Person)
	err = hltool.JSONBytesToStruct(b, p)
	if err != nil {
		log.Fatalln(err)
	}
	fmt.Println(p)
}
```

### struct序列化成二进制文件和反序列化

二进制文件可以存储到磁盘上，再次利用

```go
package main

import (
	"fmt"
	"log"

	"github.com/chanyipiaomiao/hltool"
)

// Person 人
type Person struct {
	Name    string 
	Age     int    
	Address struct {
		Country  string 
		Province string 
		City     string 
	} 
}

func main() {

	p := &Person{
		Name: "张三",
		Age:  20,
	}

	p.Address.Country = "China"
	p.Address.Province = "Shanghai"
	p.Address.City = "Shanghai"

	fmt.Println("序列化成二进制文件之前")
	fmt.Println(p)

	// 序列化成二级制文件，可以存储到磁盘上
	err := hltool.StructToBinFile(p, "/tmp/p.bin")
	if err != nil {
		log.Fatalln(err)
	}

	// 反序列化
	p2 := new(Person)
	err = hltool.BinFileToStruct("/tmp/p.bin", p2)
	if err != nil {
		log.Fatalln(err)
	}

	fmt.Println("从二进制文件中转换之后")
	fmt.Println(p2)

}

```

[返回到目录](#功能列表)

### struct序列化成byte数组和反序列化

struct序列化成byte数组，可以存储到数据库中,再次利用

```go
package main

import (
	"fmt"
	"log"

	"github.com/chanyipiaomiao/hltool"
)

// Person 人
type Person struct {
	Name    string `json:"Name"`
	Age     int    `json:"Age"`
	Address struct {
		Country  string `json:"Country"`
		Province string `json:"Province"`
		City     string `json:"City"`
	} `json:"Address"`
}

func main() {

	p := &Person{
		Name: "张三",
		Age:  20,
	}

	p.Address.Country = "China"
	p.Address.Province = "Shanghai"
	p.Address.City = "Shanghai"

	fmt.Println("struct序列化成[]byte")

	// struct序列化成[]byte，可以存储到数据库
	b, err := hltool.StructToBytes(p)
	if err != nil {
		log.Fatalln(err)
	}
	fmt.Println(p)
	fmt.Println(b)

	// []byte反序列化成struct 和序列化之前的结构体结构必须要一样
	fmt.Println("[]byte反序列化成struct")
	p2 := new(Person)
	err = hltool.BytesToStruct(b, p2)
	if err != nil {
		log.Fatalln(err)
	}
	fmt.Println(p2)

}

```

[返回到目录](#功能列表)