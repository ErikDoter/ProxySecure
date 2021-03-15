package main

import (
	"fmt"
	"net/http"
	"strings"
)

func main() {
	server := &http.Server{
		Addr: ":8887",
		Handler: http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			buf := strings.Split(r.URL.RawQuery,"&")
			var res []string
			for _, val := range buf {
				piece := strings.Split(val, "=")
				for _, val2 := range piece {
					res = append(res, val2)
				}
			}
			var nameFamilia [2]string
			var i = 0
			for index, value := range res {
				if index % 2 == 1 {
					nameFamilia[i] = value
					i++
				}
			}
			result := "Ваше имя: " + nameFamilia[0] + "Ваша фамилия" + nameFamilia[1]
			w.Write([]byte(result))
		}),
	}
	fmt.Println("listen 8887")
	server.ListenAndServe()
}