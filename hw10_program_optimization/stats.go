package hw10programoptimization

import (
	"bufio"
	"io"
	"strings"

	"github.com/tidwall/gjson"
)

//go:generate easyjson -all user.go
type User struct {
	Email string
}

type DomainStat map[string]int

func GetDomainStat(r io.Reader, domain string) (DomainStat, error) {
	return countDomainInUserData(r, domain)
}

func countDomainInUserData(r io.Reader, zone string) (DomainStat, error) {
	result := make(DomainStat)
	scanner := bufio.NewScanner(r)
	for scanner.Scan() {
		email := ExtractEmail(string(scanner.Bytes()))
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

func ExtractEmail(json string) string {
	return gjson.Get(json, "Email").String()
}
