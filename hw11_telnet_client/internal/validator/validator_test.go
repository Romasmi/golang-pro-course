package validator

import (
	"testing"
)

func TestValidateHost(t *testing.T) {
	tests := []struct {
		name    string
		host    string
		wantErr bool
	}{
		{name: "valid IPv4", host: "192.168.1.1", wantErr: false},
		{name: "valid IPv4 loopback", host: "127.0.0.1", wantErr: false},
		{name: "valid IPv4 zeros", host: "0.0.0.0", wantErr: false},

		{name: "valid IPv6", host: "2001:db8::1", wantErr: false},
		{name: "valid IPv6 loopback", host: "::1", wantErr: false},
		{name: "valid IPv6 full", host: "fe80::1ff:fe23:4567:890a", wantErr: false},

		{name: "simple domain", host: "example.com", wantErr: false},
		{name: "subdomain", host: "sub.example.com", wantErr: false},
		{name: "deep subdomain", host: "a.b.c.example.com", wantErr: false},
		{name: "hyphenated domain", host: "my-host.example.com", wantErr: false},
		{name: "numeric subdomain", host: "host1.example.com", wantErr: false},

		{name: "empty host", host: "", wantErr: true},

		{name: "localhost bare", host: "localhost", wantErr: false},

		{name: "starts with hyphen", host: "-example.com", wantErr: true},
		{name: "ends with hyphen", host: "example-.com", wantErr: true},
		{name: "double dot", host: "example..com", wantErr: true},
		{name: "trailing dot", host: "example.com.", wantErr: true},
		{name: "only zone", host: ".com", wantErr: true},
		{name: "invalid IP octet", host: "999.999.999.999", wantErr: true},
		{name: "spaces in host", host: "exa mple.com", wantErr: true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateHost(tt.host)
			if (err != nil) != tt.wantErr {
				t.Errorf("ValidateHost(%q) error = %v, wantErr = %v", tt.host, err, tt.wantErr)
			}
		})
	}
}

func TestValidatePort(t *testing.T) {
	tests := []struct {
		name    string
		port    int
		wantErr bool
	}{
		{name: "min port", port: 1, wantErr: false},
		{name: "max port", port: 65535, wantErr: false},
		{name: "common HTTP", port: 80, wantErr: false},
		{name: "common HTTPS", port: 443, wantErr: false},
		{name: "common app port", port: 8080, wantErr: false},

		{name: "port zero", port: 0, wantErr: true},
		{name: "port above max", port: 65536, wantErr: true},
		{name: "negative port", port: -1, wantErr: true},
		{name: "large number", port: 99999, wantErr: true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidatePort(tt.port)
			if (err != nil) != tt.wantErr {
				t.Errorf("ValidatePort(%q) error = %v, wantErr = %v", tt.port, err, tt.wantErr)
			}
		})
	}
}
