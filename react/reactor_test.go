package react

import (
	"net/http"
	"testing"
)

func BenchmarkGet(b *testing.B) {
	b.StopTimer()
	b.StartTimer()
	for i := 0; i < b.N; i++ {
		c := &http.Client{}
		c.Get("http://localhost:8080/index/3")
	}
}
