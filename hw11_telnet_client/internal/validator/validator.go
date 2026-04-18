package validator

import (
	"fmt"
	"net"
	"regexp"
)

var domainRegex = regexp.MustCompile(`^([a-zA-Z0-9]([a-zA-Z0-9\-]{0,61}[a-zA-Z0-9])?\.)+[a-zA-Z]{2,}$`)

func ValidateHost(host string) error {
	if host == "" {
		return fmt.Errorf("host cannot be empty")
	}
	if net.ParseIP(host) != nil {
		return nil
	}
	if host == "localhost" || domainRegex.MatchString(host) {
		return nil
	}
	return fmt.Errorf("invalid host %q: must be a valid IP address or domain name", host)
}

func ValidatePort(port int) error {
	if port < 1 || port > 65535 {
		return fmt.Errorf("invalid port %d: must be between 1 and 65535", port)
	}
	return nil
}
