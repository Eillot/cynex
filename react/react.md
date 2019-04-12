## HTTP Server 和 Router构建

#### 提供HTTP-Server启动工具集；提供快速构建View层开发的方式。

#### 用法

##### 1、创建程序启动文件main.go，在main方法中启动HttpServer
```
package main

import (
	"cynex/react"
)

func main() {
	react.Server.Start()
}

```

##### 2、创建控制器文件controller/user.go；在文件中定义控制器类型User，在init方法中绑定路径处理方法
```
package controller

import (
	"cynex/log"
	"cynex/react"
	"net/http"
)

type User struct {
}

func init() {
	react.BindGet("/index/{flag}", new(User), "Index")
}

func (u *User) Index(w http.ResponseWriter, r *http.Request) {
	log.Debug("Index 方法已经执行")
}


```

##### 3、使用默认配置，执行main方法。

