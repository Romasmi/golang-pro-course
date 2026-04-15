package hw10programoptimization

import (
	"bufio"
	"io"
	"strings"
)

type User struct {
	ID       int
	Name     string
	Username string
	Email    string
	Phone    string
	Password string
	Address  string
}

type DomainStat map[string]int

func GetDomainStat(r io.Reader, domain string) (DomainStat, error) {
	return countDomainInUserData(r, domain)
}

func countDomainInUserData(r io.Reader, zone string) (DomainStat, error) {
	result := make(DomainStat)
	scanner := bufio.NewScanner(r)
	for scanner.Scan() {
		email := extractEmail(string(scanner.Bytes()))
		if domain, ok := getDomainInZone(email, zone); ok {
			result[domain]++
		}
	}
	return result, nil
}

func getDomainInZone(email, target string) (string, bool) {
	i := strings.IndexByte(email, '@')
	if i == -1 || i+1 >= len(email) {
		return "", false
	}

	domain := email[i+1:]

	if domain == target {
		return strings.ToLower(domain), true
	}

	for j := 0; j < len(domain); j++ {
		if domain[j] == '.' && j+1 < len(domain) {
			if domain[j+1:] == target {
				return strings.ToLower(domain), true
			}
		}
	}

	return "", false
}

func extractEmail(s string) string {
	key := `"Email":"`
	for i := 0; i <= len(s)-len(key); i++ {
		if s[i] == '"' && s[i:i+len(key)] == key {
			start := i + len(key)

			for j := start; j < len(s); j++ {
				if s[j] == '"' {
					return s[start:j]
				}
			}
			return ""
		}
	}
	return ""
}
