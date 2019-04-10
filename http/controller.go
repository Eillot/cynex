package http

import (
	"errors"
	"net/http"
	"strconv"
	"strings"
	"time"
)

type Controller struct {
}

// Cookie 用于读取Cookie中的值
func (c *Controller) Cookie(name string) (string, error) {
	cookie, err := defaultReactor.Request.Cookie(name)
	if err == nil {
		return cookie.Value, nil
	}
	return "", err
}

// SetCookie 用于设置Cookie值
// 参数 name: 名称；value: 值；expireAfter: N秒后过期，单位秒
func (c *Controller) SetCookie(name, value string, expireAfter int) {
	d, _ := time.ParseDuration(strconv.Itoa(expireAfter) + "s")
	cookie := &http.Cookie{
		Name:     name,
		Value:    value,
		Expires:  time.Now().Add(d),
		MaxAge:   expireAfter,
		HttpOnly: false,
		SameSite: http.SameSiteStrictMode,
	}
	http.SetCookie(defaultReactor.ResponseWriter, cookie)
}

// 读取参数值
func (c *Controller) FormValue(name string) (string, error) {
	if val := defaultReactor.Request.FormValue(name); strings.TrimSpace(val) != "" {
		return strings.TrimSpace(val), nil
	}
	return "", errors.New("empty")
}

// Request 用于提供当前请求*http.Request的指针
func (c *Controller) Request() *http.Request {
	return defaultReactor.Request
}

// ResponseWriter 用于提供当前请求http.ResponseWriter
func (c *Controller) ResponseWriter() http.ResponseWriter {
	return defaultReactor.ResponseWriter
}
