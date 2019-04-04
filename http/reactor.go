package http

import "net/http"

type Reactor struct {
	r *http.Request
	w http.ResponseWriter
}
