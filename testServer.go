package main

import (
	"fmt"
	"net/http"
)

func main() {
	server := &http.Server{
		Addr: ":8887",
		Handler: http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			name := r.URL.Query().Get("name")
			familia := r.URL.Query().Get("familia")
			result := "Ваше имя: " + name + " Ваша фамилия: " + familia
			w.Write([]byte(result))
			w.WriteHeader(200)
		}),
	}
	fmt.Println("listen 8887")
	server.ListenAndServe()
}
