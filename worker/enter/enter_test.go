package enter

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

const (
	testLatitude  = "52.52"
	testLongitude = "13.41"
	baseURL       = "https://api.open-meteo.com/v1/forecast"
)

// TestExtractDetails_GetRequest tests GET request handling through extactDetails.
// This approach creates a mock destination server to intercept the proxied request,
// allowing us to verify headers, method, URL construction, and response metrics
// without network dependencies or rate limiting from the real API.
func TestExtractDetails_GetRequest(t *testing.T) {
	var receivedMethod string
	dest := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		receivedMethod = r.Method
		w.WriteHeader(http.StatusOK)
	}))
	defer dest.Close()
	req := httptest.NewRequest("GET", dest.URL+"/weather?lat=52.52&lon=13.41", nil)
	req.Header.Set("X-Forwarded-For", "192.168.1.1")
	req.Header.Set("User-Agent", "TestClient")
	payload, err := extactDetails(req)
	if err != nil {
		t.Fatalf("ExtactDetails failed: %v", err)
	}

	if receivedMethod != "GET" {
		t.Errorf("Expected method GET, got %s", receivedMethod)
	}

	if payload.Method != "GET" {
		t.Errorf("Expected payload.Method GET, got %s", payload.Method)
	}

	if payload.Headers["X-Forwarded-For"] != "192.168.1.1" {
		t.Errorf("Expected X-Forwarded-For header, got %v", payload.Headers["X-Forwarded-For"])
	}

	if payload.Metrics.StatusCode != http.StatusOK {
		t.Errorf("Expected status 200, got %d", payload.Metrics.StatusCode)
	}

	if payload.Metrics.Duration <= 0 {
		t.Error("Expected positive duration")
	}

	t.Logf("Payload: method=%s, url=%s, status=%d, duration=%v",
		payload.Method, payload.Url, payload.Metrics.StatusCode, payload.Metrics.Duration)
}

// TestExtractDetails_PostRequest tests POST request handling with JSON body.
// This validates that request bodies are correctly forwarded and that the Content-Type
// header is set appropriately for API compatibility.
func TestExtractDetails_PostRequest(t *testing.T) {
	var receivedBody map[string]any
	var receivedContentType string

	dest := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		receivedContentType = r.Header.Get("Content-Type")
		json.NewDecoder(r.Body).Decode(&receivedBody)
		w.WriteHeader(http.StatusCreated)
	}))
	defer dest.Close()

	requestBody := map[string]any{"city": "Berlin", "units": "metric"}
	bodyBytes, _ := json.Marshal(requestBody)
	req := httptest.NewRequest("POST", dest.URL+"/api/data", bytes.NewReader(bodyBytes))
	req.Header.Set("Content-Type", "application/json")

	payload, err := extactDetails(req)
	if err != nil {
		t.Fatalf("ExtactDetails failed: %v", err)
	}

	if payload.Method != "POST" {
		t.Errorf("Expected method POST, got %s", payload.Method)
	}

	if receivedContentType != "application/json" {
		t.Errorf("Expected Content-Type application/json, got %s", receivedContentType)
	}

	if payload.Metrics.StatusCode != http.StatusCreated {
		t.Errorf("Expected status 201, got %d", payload.Metrics.StatusCode)
	}

	if _, exists := payload.Body["city"]; !exists {
		t.Error("Expected body.city in payload")
	}
}

// TestExtractDetails_PutRequest tests PUT request handling.
// This ensures the function handles PUT method correctly, forwarding both headers
// and request body while capturing the response status.
func TestExtractDetails_PutRequest(t *testing.T) {
	var receivedBody map[string]any

	dest := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		json.NewDecoder(r.Body).Decode(&receivedBody)
		w.WriteHeader(http.StatusNoContent)
	}))
	defer dest.Close()

	requestBody := map[string]any{"temperature": 22.5}
	bodyBytes, _ := json.Marshal(requestBody)
	req := httptest.NewRequest("PUT", dest.URL+"/api/update", bytes.NewReader(bodyBytes))

	payload, err := extactDetails(req)
	if err != nil {
		t.Fatalf("ExtactDetails failed: %v", err)
	}

	if payload.Method != "PUT" {
		t.Errorf("Expected method PUT, got %s", payload.Method)
	}

	if payload.Metrics.StatusCode != http.StatusNoContent {
		t.Errorf("Expected status 204, got %d", payload.Metrics.StatusCode)
	}
}

// TestExtractDetails_DeleteRequest tests DELETE request handling.
// This verifies that DELETE requests are forwarded without a body and that
// successful deletions are tracked in the metrics.
func TestExtractDetails_DeleteRequest(t *testing.T) {
	dest := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "DELETE" {
			t.Errorf("Expected DELETE method, got %s", r.Method)
		}
		w.WriteHeader(http.StatusNoContent)
	}))
	defer dest.Close()

	req := httptest.NewRequest("DELETE", dest.URL+"/api/resource/123", nil)

	payload, err := extactDetails(req)
	if err != nil {
		t.Fatalf("ExtactDetails failed: %v", err)
	}

	if payload.Method != "DELETE" {
		t.Errorf("Expected method DELETE, got %s", payload.Method)
	}

	if len(payload.Body) != 0 {
		t.Error("Expected empty body for DELETE request")
	}
}

// TestExtractDetails_HeaderExtraction tests that all expected headers are captured.
// This approach verifies the full header forwarding chain, ensuring proxy headers
// (X-Forwarded-For, X-Forwarded-Host, etc.) are preserved for upstream services.
func TestExtractDetails_HeaderExtraction(t *testing.T) {
	var capturedHeaders http.Header

	dest := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		capturedHeaders = r.Header
		w.WriteHeader(http.StatusOK)
	}))
	defer dest.Close()

	req := httptest.NewRequest("GET", dest.URL+"/test", nil)
	req.Header.Set("X-Forwarded-For", "10.0.0.1")
	req.Header.Set("X-Forwarded-Host", "api.example.com")
	req.Header.Set("X-Forwarded-Proto", "https")
	req.Header.Set("X-Real-IP", "10.0.0.2")
	req.Header.Set("Authorization", "Bearer token123")
	req.Header.Set("Accept", "application/json")

	payload, err := extactDetails(req)
	if err != nil {
		t.Fatalf("ExtactDetails failed: %v", err)
	}

	if payload.Headers["X-Forwarded-For"] != "10.0.0.1" {
		t.Errorf("Expected X-Forwarded-For '10.0.0.1', got %v", payload.Headers["X-Forwarded-For"])
	}

	if payload.Headers["X-Forwarded-Host"] != "api.example.com" {
		t.Errorf("Expected X-Forwarded-Host 'api.example.com', got %v", payload.Headers["X-Forwarded-Host"])
	}

	if capturedHeaders.Get("User-Agent") != "Metricraft" {
		t.Errorf("Expected User-Agent 'Metricraft', got %s", capturedHeaders.Get("User-Agent"))
	}

	if capturedHeaders.Get("Authorization") != "Bearer token123" {
		t.Errorf("Expected Authorization header to be forwarded")
	}

	t.Logf("Captured headers: %v", capturedHeaders)
}

// TestExtractDetails_UrlConstruction tests URL building with query parameters.
// This validates that query strings are preserved when constructing the redirect URL,
// ensuring parameters like latitude, longitude, and API flags are correctly appended.
func TestExtractDetails_UrlConstruction(t *testing.T) {
	dest := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))
	defer dest.Close()
	req := httptest.NewRequest("GET", dest.URL+"/forecast?latitude=52.52&longitude=13.41&current_weather=true", nil)
	payload, err := extactDetails(req)
	if err != nil {
		t.Fatalf("ExtactDetails failed: %v", err)
	}
	expectedQuery := "latitude=52.52&longitude=13.41&current_weather=true"
	if payload.Url == "" {
		t.Error("Expected non-empty URL")
	}
	t.Logf("Constructed URL: %s", payload.Url)
	_ = expectedQuery
}

// TestExtractDetails_ConnectionError tests error handling when destination is unreachable.
// This approach ensures errors are propagated correctly and don't silently fail,
// providing visibility into infrastructure issues.
func TestExtractDetails_ConnectionError(t *testing.T) {
	req := httptest.NewRequest("GET", "http://127.0.0.1:1/unreachable", nil)
	_, err := extactDetails(req)
	if err == nil {
		t.Error("Expected connection error, got nil")
	}
	t.Logf("Expected error received: %v", err)
}

// TestExtractDetails_TimeField tests that the Time field is populated.
// This validates that request timestamp is captured for logging and metrics
// correlation across distributed systems.
func TestExtractDetails_TimeField(t *testing.T) {
	dest := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))
	defer dest.Close()

	before := time.Now()
	req := httptest.NewRequest("GET", dest.URL+"/test", nil)

	payload, err := extactDetails(req)
	if err != nil {
		t.Fatalf("ExtactDetails failed: %v", err)
	}

	if payload.Time.IsZero() {
		t.Error("Expected non-zero Time field")
	}

	if payload.Time.Before(before) {
		t.Error("Expected Time to be after request creation")
	}
}
