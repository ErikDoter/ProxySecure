package server

import (
	"bytes"
	"fmt"
	"github.com/jackc/pgx"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"strings"
)

type requests struct {
	Id int64
	Method string
	Url string
	Headers string
	Body string
}

type AllRequests []requests

type writer struct {
	bytes.Buffer
}

func (b *writer) Close() error {
	b.Buffer.Reset()
	return nil
}

func ConnectDb() *pgx.ConnPool {
	config := pgx.ConnPoolConfig{
		ConnConfig:     pgx.ConnConfig{
			Host:                 "localhost",
			Port:                 5432,
			Database:             "proxy",
			User:                 "proxy",
			Password:             "proxy",
			TLSConfig:            nil,
			UseFallbackTLS:       false,
			FallbackTLSConfig:    nil,
			Logger:               nil,
			LogLevel:             0,
			Dial:                 nil,
			RuntimeParams:        nil,
			OnNotice:             nil,
			CustomConnInfo:       nil,
			CustomCancel:         nil,
			PreferSimpleProtocol: false,
			TargetSessionAttrs:   "",
		},
		MaxConnections: 100,
		AfterConnect:   nil,
		AcquireTimeout: 0,
	}
	connPool, err := pgx.NewConnPool(config)
	if err != nil {
		log.Fatal(err)
	}
	return connPool
}

func SaveRequest(db *pgx.ConnPool, r *http.Request) error {
	sql := `INSERT INTO requests VALUES(default,$1,$2,$3,$4)`
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		fmt.Println(err)
		return err
	}
	headers := ""
	for key, val := range r.Header {
		headers += key + ": " + val[0] + "\n"
	}
	queryResult, err := db.Exec(sql,
		r.Method, r.URL.String(), headers, string(body))
	affected := queryResult.RowsAffected()
	if (affected != 1) || (err != nil) {
		fmt.Print(err)
		return err
	}
	return nil
}

func GetRequest(id int, db *pgx.ConnPool) http.Request {
	var result http.Request
	var request requests
	err := db.QueryRow("select id, method, url, headers, body from requests where id = $1", id).Scan(&request.Id, &request.Method, &request.Url, &request.Headers, &request.Body)
	if err != nil {
		fmt.Println(err)
		return http.Request{}
	}
	result.Method = request.Method
	result.URL, err = url.Parse(request.Url)
	if err != nil {
		fmt.Println(err)
		return http.Request{}
	}
	var bodyWriter io.ReadWriteCloser
	bodyWriter = &writer{}
	_, err = bodyWriter.Write([]byte(request.Body))
	if err != nil {
		fmt.Println(err)
		return http.Request{}
	}
	result.Body = bodyWriter
	headMap := make(map[string][]string)
	for _, val := range strings.Split(request.Headers, "\n") {
		if val != "" {
			buf := strings.Split(val, ":")
			headMap[buf[0]] = []string{buf[1]}
		}
	}
	result.Header = headMap
	return result
}

func GetAllRequests(db *pgx.ConnPool) AllRequests{
	var allReq AllRequests
	rows, err := db.Query("select * from requests")
	if err != nil {
		return allReq
	}
	defer rows.Close()
	for rows.Next() {
		request := requests{}
		err := rows.Scan(&request.Id, &request.Method, &request.Url, &request.Headers, &request.Body)
		if err != nil {
			return allReq
		}
		allReq = append(allReq, request)
	}
	return allReq
}