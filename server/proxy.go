package server

import (
	"encoding/json"
	"github.com/jackc/pgx"
	"io"
	"net"
	"net/http"
	"strconv"
	"strings"
	"time"
)

func HandleTunneling(w http.ResponseWriter, r *http.Request) {
	dest_conn, err := net.DialTimeout("tcp", r.Host, 10*time.Second)
	if err != nil {
		http.Error(w, err.Error(), http.StatusServiceUnavailable)
		return
	}
	w.WriteHeader(http.StatusOK)
	hijacker, ok := w.(http.Hijacker)
	if !ok {
		http.Error(w, "Hijacking not supported", http.StatusInternalServerError)
		return
	}
	client_conn, _, err := hijacker.Hijack()
	if err != nil {
		http.Error(w, err.Error(), http.StatusServiceUnavailable)
	}
	go transfer(dest_conn, client_conn)
	go transfer(client_conn, dest_conn)
}

func transfer(destination io.WriteCloser, source io.ReadCloser) {
	defer destination.Close()
	defer source.Close()
	io.Copy(destination, source)
}

func HandleHTTP(w http.ResponseWriter, req *http.Request) {
	req.Header.Del("Proxy-Connection")
	resp, err := http.DefaultTransport.RoundTrip(req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusServiceUnavailable)
		return
	}
	defer resp.Body.Close()
	copyHeader(w.Header(), resp.Header)
	w.WriteHeader(resp.StatusCode)
	io.Copy(w, resp.Body)
}

func copyHeader(dst, src http.Header) {
	for k, vv := range src {
		for _, v := range vv {
			dst.Add(k, v)
		}
	}
}

func RepeatRequest(url string, w http.ResponseWriter, r *http.Request, db *pgx.ConnPool) {
	buffer := strings.Split(url, "/")
	id, err := strconv.Atoi(buffer[2])
	if err != nil {
		return
	}
	request := GetRequest(id, db)
	r = &request
	http.Redirect(w, &request, request.URL.String(), 301)

	return
}

func RequestList(w http.ResponseWriter, r *http.Request, db *pgx.ConnPool) {
	result := GetAllRequests(db)
	answer, err := json.Marshal(result)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(200)
	w.Write(answer)
}