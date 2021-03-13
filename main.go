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
		Addr: ":8888",
		Handler: http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if r.Method == http.MethodConnect {
				server.HandleTunneling(w, r)
			} else  {
				repeaterPattern := `^/request/[0-9]+$`
				if match, _ := regexp.Match(repeaterPattern, []byte(r.URL.String())); match {
					fmt.Println(r.URL, "  ", match)
				} else {
					server.SaveRequest(db, r)
					server.HandleHTTP(w, r)
				}
			}
		}),
		TLSNextProto: make(map[string]func(*http.Server, *tls.Conn, http.Handler)),
	}
	fmt.Println("listen 8888")
	server.ListenAndServe()
}
