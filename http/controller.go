package http

import (
	"net/http"
)

var defaultController *Controller

type Controller struct {
	handler *reactor
}

func init() {
	defaultController = &Controller{
		handler: handler,
	}
}

func (c *Controller) GetCookie(name string) string {
	if cookie, err := c.handler.Request.Cookie(name); err == nil {
		return cookie.Value
	}
	return ""
}

func (c *Controller) SetCookie(name, value string, expireAfter int) {
	cookie := &http.Cookie{
		Name:     name,
		Value:    value,
		MaxAge:   expireAfter,
		HttpOnly: false,
	}
	c.handler.ResponseWriter.Header().Set("Set-Cookie", cookie.String())
}
