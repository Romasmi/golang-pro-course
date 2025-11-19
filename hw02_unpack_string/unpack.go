package hw02unpackstring

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
	"unicode"
)

var ErrInvalidString = errors.New("invalid string")

func Unpack(packed string) (string, error) {
	var unpacked strings.Builder
	var prevR rune
	isPrevEscaped := false
	for _, r := range packed {
		switch {
		case unicode.IsDigit(r):
			if !isPrevEscaped && prevR == '\\' {
				prevR = r
				isPrevEscaped = true
				continue
			}

			if !isPrevEscaped && unicode.IsDigit(prevR) {
				return "", fmt.Errorf("%w: numbers are not accepted", ErrInvalidString)
			}
			if prevR == 0 {
				return "", fmt.Errorf("%w: digit can not be before a substring", ErrInvalidString)
			}

			digit, _ := strconv.Atoi(string(r))
			unpacked.WriteString(strings.Repeat(string(prevR), digit))
			prevR = 0
		case prevR == '\\' && !isPrevEscaped:
			prevR = r
			isPrevEscaped = true
			continue
		default:
			if prevR != 0 {
				unpacked.WriteRune(prevR)
			}
			prevR = r
		}
		isPrevEscaped = false
	}
	if prevR != 0 {
		unpacked.WriteRune(prevR)
	}
	return unpacked.String(), nil
}
