package main

import (
	"crypto/tls"
	"fmt"
	"github.com/ErikDoter/ProxySecure/server"
	"net/http"
	"regexp"
)

func main() {
	db := server.ConnectDb()
	defer db.Close()
	server := &http.Server{
		Addr: ":8080",
		Handler: http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if r.Method == http.MethodConnect {
				server.HandleTunneling(w, r)
			} else  {
				repeaterPattern := `^/request/[0-9]+$`
				requestsPattern := `^/requests$`
				xssPatern := `^/xss/[0-9]+`
				if match, _ := regexp.Match(repeaterPattern, []byte(r.URL.String())); match {
					server.RepeatRequest(r.URL.String(), w, r, db)
				} else if match, _ := regexp.Match(requestsPattern, []byte(r.URL.String())); match {
					server.RequestList(w, r, db)
				} else if match, _ := regexp.Match(xssPatern, []byte(r.URL.String())); match {
					server.CheckXSS(w, r, r.URL.String(), db)
				} else {
					server.SaveRequest(db, r)
					server.HandleHTTP(w, r)
				}
			}
		}),
		TLSNextProto: make(map[string]func(*http.Server, *tls.Conn, http.Handler)),
	}
	fmt.Println("listen 8888")
	fmt.Println("/requests - Список всех запросов")
	fmt.Println("/request/{id} - Повторить запрос с этим id")
	fmt.Println("/xss/{id} - Проверить на уязвимость xss запрос с этим id")
	server.ListenAndServe()
}
