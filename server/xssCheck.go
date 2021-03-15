package server

import (
	"fmt"
	"github.com/jackc/pgx"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"
)

func CheckXSS(w http.ResponseWriter, r *http.Request, url string, db *pgx.ConnPool) {
	attackVector := "vulnerable'\"><img%20src%20on%20error=alert()>"
	var flag bool
	buffer := strings.Split(url, "/")
	id, err := strconv.Atoi(buffer[2])
	if err != nil {
		return
	}
	request := GetRequest(id, db)
	if request.Method == "" {
		fmt.Println("request doesn't exist")
		return
	}
	oldQuery := request.URL.RawQuery
	for key, value := range request.URL.Query() {
		request.URL.RawQuery = strings.ReplaceAll(request.URL.RawQuery, value[0], attackVector)
		resp, err := http.DefaultTransport.RoundTrip(&request)
		if err != nil {
			fmt.Println("error with round trip")
			return
		}

		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			fmt.Println("error with Read all")
			return
		}
		if strings.Contains(string(body), "vulnerable'\"><img src on error=alert()>") {
			w.Write([]byte(key + " - уязвимый параметр\n"))
			fmt.Println(key + " - уязвимый параметр")
			flag = true
		} else {
			w.Write([]byte(key + " - уязвимости не найдено\n"))
			fmt.Println(key + " - уязвимости не найдено")
			flag = true
		}
		request.URL.RawQuery = oldQuery
		resp.Body.Close()
	}
	if !flag {
		w.Write([]byte("Параметров не найдено"))
		fmt.Println("Параметров не найдено")
	}
}