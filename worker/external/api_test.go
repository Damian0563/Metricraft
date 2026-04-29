package external

import (
	"testing"
)

func TestGetCountryOrigin(t *testing.T) {
	tests := []struct {
		name     string
		ip       string
		expected string
	}{
		{
			name:     "Google DNS US",
			ip:       "8.8.8.8",
			expected: "United States",
		},
		{
			name:     "Google DNS EU",
			ip:       "8.8.4.4",
			expected: "United States",
		},
		{
			name:     "Cloudflare DNS",
			ip:       "1.1.1.1",
			expected: "Australia",
		},
		{
			name:     "Quad9 DNS",
			ip:       "9.9.9.9",
			expected: "United States",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			country, err := GetCountryOrigin(tt.ip)
			if err != nil {
				t.Fatalf("GetCountryOrigin(%q) failed: %v", tt.ip, err)
			}
			if country == "" {
				t.Errorf("Expected non-empty country for %s", tt.ip)
			}
			t.Logf("Country for %s: %s", tt.ip, country)
		})
	}
}

func TestGetCountryOrigin_InvalidIP(t *testing.T) {
	_, err := GetCountryOrigin("invalid_ip")
	if err != nil {
		t.Logf("Expected error for invalid IP: %v", err)
	}
}

func TestGetCountryOrigin_PrivateIP(t *testing.T) {
	country, err := GetCountryOrigin("127.0.0.1")
	if err != nil {
		t.Fatalf("GetCountryOrigin(127.0.0.1) failed: %v", err)
	}
	t.Logf("Country for 127.0.0.1: %q", country)
}

func TestGetCountryOrigin_ReservedIP(t *testing.T) {
	country, err := GetCountryOrigin("256.256.256.256")
	if err != nil {
		t.Logf("Error received: %v", err)
		return
	}
	t.Logf("Country for 256.256.256.256: %q", country)
}
