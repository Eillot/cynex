package http

import (
	"net/http"
)

// 被引入时完成初始化
var server = &http.Server{}

/**
 * 启动服务监听
 */
func Run() {
	server.ListenAndServe()
}

/**
 * 处理指定路径的请求
 * @path 请求路径
 * @reactor 请求处理组件块
 * @method  处理方法（rest方法使用冒号间隔，例如List:GET;多个方法使用逗号分隔(List,First)）
 */
func Accept(path string, reactor interface{}, method ...string) {
	a := Reactor(reactor)
	server.Handler.ServeHTTP(a.w, a.r)
}
