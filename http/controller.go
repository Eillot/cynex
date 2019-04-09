package http

import (
	"net/http"
	"strconv"
	"time"
)

type Controller struct {
}

// GetCookie 用于读取Cookie中的值
func (c *Controller) GetCookie(name string) (string, error) {
	cookie, err := handler.Request.Cookie(name)
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
	http.SetCookie(handler.ResponseWriter, cookie)
}
