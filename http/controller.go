package http

import (
	"net/http"
)

type Controller struct {
}

func (c *Controller) GetCookie(name string) string {
	if cookie, err := handler.Request.Cookie(name); err == nil {
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
		SameSite: http.SameSiteStrictMode,
	}
	http.SetCookie(handler.ResponseWriter, cookie)
}
