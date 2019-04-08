package http

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
)

type User struct {
}

func (u *User) Index() {
	fmt.Println("INdex is Executing!")
}

func (u *User) Hello() {
	fmt.Println("Hello is Executing!")
}

func TestGet(t *testing.T) {
	Get("/index/[*nd*]/{index}", &User{}, "Index")
	Post("/hello", &User{}, "Hello")
	server.Run()
}

func BenchmarkGet(b *testing.B) {
	for i := 0; i < b.N; i++ {
		r := httptest.NewRecorder()
		c := &http.Client{}
		c.Get("http://localhost:8080/index/index/3")
		r.Result()
	}
}
