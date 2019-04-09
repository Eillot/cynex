package http

import "net/http"

func Start() {
	http.ListenAndServe(":8080", reactor)
}
