package enter

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
	"time"
)

var client = &http.Client{Timeout: 30 * time.Second}

func Enter(w http.ResponseWriter, r *http.Request) {
	payload, err := extactDetails(r)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "Error: %s", err)
		return
	}
	Leave(payload)
}

func extactDetails(r *http.Request) (Payload, error) {
	started := time.Now()
	var headers map[string]string = make(map[string]string)
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
	destPort := os.Getenv("DEST_PORT")
	if destPort != "" && !strings.Contains(host, ":") {
		host = fmt.Sprintf("%s:%s", host, destPort)
	}
	redirect := fmt.Sprintf("%s://%s%s", scheme, host, r.URL.Path)
	if r.URL.RawQuery != "" {
		redirect += "?" + r.URL.RawQuery
	}
	var metrics ResponseMetrics
	start := time.Now()
	switch method {
	case "GET":
		req, err := http.NewRequest("GET", redirect, nil)
		if err != nil || req == nil {
			return Payload{}, fmt.Errorf("failed to create request: %v", err)
		}
		if r.Header != nil {
			for k, v := range r.Header {
				req.Header.Set(k, v[0])
			}
		}
		req.Header.Set("User-Agent", "Metricraft")
		resp, err := client.Do(req)
		if err == nil {
			metrics.StatusCode = resp.StatusCode
			resp.Body.Close()
		} else {
			return Payload{}, err
		}
	case "POST":
		req, _ := http.NewRequest("POST", redirect, bytes.NewBuffer(data))
		if r.Header != nil {
			for k, v := range r.Header {
				req.Header.Set(k, v[0])
			}
		}
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("User-Agent", "Metricraft")
		resp, err := client.Do(req)
		if err == nil {
			metrics.StatusCode = resp.StatusCode
			resp.Body.Close()
		} else {
			return Payload{}, err
		}
	case "PUT":
		req, _ := http.NewRequest("PUT", redirect, bytes.NewBuffer(data))
		if r.Header != nil {
			for k, v := range r.Header {
				req.Header.Set(k, v[0])
			}
		}
		req.Header.Set("User-Agent", "Metricraft")
		req.Header.Set("Content-Type", "application/json")
		resp, err := client.Do(req)
		if err == nil {
			metrics.StatusCode = resp.StatusCode
			resp.Body.Close()
		} else {
			return Payload{}, err
		}
	case "DELETE":
		req, _ := http.NewRequest("DELETE", redirect, nil)
		if r.Header != nil {
			for k, v := range r.Header {
				req.Header.Set(k, v[0])
			}
		}
		req.Header.Set("User-Agent", "Metricraft")
		resp, err := client.Do(req)
		if err == nil {
			metrics.StatusCode = resp.StatusCode
			resp.Body.Close()
		} else {
			return Payload{}, err
		}
	}
	metrics.Duration = time.Since(start).Milliseconds()
	return Payload{Headers: headers, Time: started, Url: redirect, Body: body, Method: method, Metrics: metrics}, nil
}
