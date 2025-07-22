package main

// 练习场代码，用来测试
import (
	"io"
	"net/http"
)

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("/ping", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("pong"))
		io.WriteString(w, "pong")
	})
	if err := http.ListenAndServe(":8080", mux); err != nil {
		panic(err)
	}

}
