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
	buf := strings.Split(request.URL.RawQuery,"&")
	var res []string
	for _, val := range buf {
		piece := strings.Split(val, "=")
		for _, val2 := range piece {
			res = append(res, val2)
		}
	}
	var answer string
	for index, value := range res {
		if index % 2 == 0 {
			answer = answer + value
			if index != len(res) - 1 {
				answer += "="
			}
		} else if index % 2 == 1 {
			answer = answer + "vulnerable'\"><img src onerror=alert()>"
			if index != len(res) - 1 {
				answer += "&"
			}
		}
	}
	request.URL.RawQuery = answer
	resp, err := http.DefaultTransport.RoundTrip(&request)
	if err != nil {
		fmt.Println("error with round trip")
		return
	}
	fmt.Println(resp)
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("error with Read all")
		return
	}
	fmt.Println(body)
}