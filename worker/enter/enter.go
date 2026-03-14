package enter

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"time"
)

var client = &http.Client{Timeout: 30 * time.Second}

func Enter(w http.ResponseWriter, r *http.Request) {
	var headers map[string]any = make(map[string]any)
	headers["X-Forwarded-For"] = r.Header.Get("X-Forwarded-For")
	headers["X-Forwarded-Host"] = r.Header.Get("X-Forwarded-Host")
	headers["X-Forwarded-Proto"] = r.Header.Get("X-Forwarded-Proto")
	headers["X-Real-IP"] = r.Header.Get("X-Real-IP")
	method := r.Method
	var body map[string]any
	data, _ := io.ReadAll(r.Body)
	if len(data) > 0 {
		json.Unmarshal(data, &body)
	}
	scheme := "http"
	if r.TLS != nil {
		scheme = "https"
	}
	host := r.Host
	redirect := fmt.Sprintf("%s://%s:%s%s", scheme, host, os.Getenv("DEST_PORT"), r.URL.Path)
	if r.URL.RawQuery != "" {
		redirect += "?" + r.URL.RawQuery
	}
	var metrics ResponseMetrics
	start := time.Now()
	switch method {
	case "GET":
		req, _ := http.NewRequest("GET", redirect, nil)
		for k, v := range r.Header {
			req.Header.Set(k, v[0])
		}
		resp, err := client.Do(req)
		if err == nil {
			metrics.StatusCode = resp.StatusCode
			resp.Body.Close()
		}
	case "POST":
		req, _ := http.NewRequest("POST", redirect, bytes.NewBuffer(data))
		for k, v := range r.Header {
			req.Header.Set(k, v[0])
		}
		req.Header.Set("Content-Type", "application/json")
		resp, err := client.Do(req)
		if err == nil {
			metrics.StatusCode = resp.StatusCode
			resp.Body.Close()
		}
	case "PUT":
		req, _ := http.NewRequest("PUT", redirect, bytes.NewBuffer(data))
		for k, v := range r.Header {
			req.Header.Set(k, v[0])
		}
		req.Header.Set("Content-Type", "application/json")
		resp, err := client.Do(req)
		if err == nil {
			metrics.StatusCode = resp.StatusCode
			resp.Body.Close()
		}
	case "DELETE":
		req, _ := http.NewRequest("DELETE", redirect, nil)
		for k, v := range r.Header {
			req.Header.Set(k, v[0])
		}
		resp, err := client.Do(req)
		if err == nil {
			metrics.StatusCode = resp.StatusCode
			resp.Body.Close()
		}
	}
	metrics.Duration = time.Since(start)
	payload := Payload{Headers: headers, Url: redirect, Body: body, Method: method, Metrics: metrics}
	Leave(payload)
}
