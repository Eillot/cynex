package react

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
)

type User struct {
	reactor
}

func (u *User) Index() {
	fmt.Println("INdex is Executing!")
}

func (u *User) Hello() {
	fmt.Println("Hello is Executing!")
}

func TestGet(t *testing.T) {
	BindGet("/index/[*nd*]/{index}", &User{}, "Index")
	BindGet("/index/rbdex/{bbb}", &User{}, "Hello")
	BindPost("/index", new(User), "Hello")
	BindDownload("/download", "1.jpg")
	BindStatic("/static", "/static")

	Server.Start()
}

func BenchmarkGet(b *testing.B) {
	for i := 0; i < b.N; i++ {
		r := httptest.NewRecorder()
		c := &http.Client{}
		c.Get("react://localhost:8080/index/index/3")
		r.Result()
	}
}
