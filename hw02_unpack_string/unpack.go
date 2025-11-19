package hw02unpackstring

import (
	"errors"
	"strconv"
	"strings"
	"unicode"
	"unicode/utf8"
)

var (
	ErrNumbersAreNotAllowed        = errors.New("numbers > 9 not allowed")
	ErrStringCanNotStartWithADigit = errors.New("string can not start with a digit")
	ErrInvalidEscaping             = errors.New("invalid escaping")
)

func Unpack(packed string) (string, error) {
	if packed == "" {
		return packed, nil
	}
	runes := []rune(packed)
	if unicode.IsDigit(runes[0]) {
		return "", ErrStringCanNotStartWithADigit
	}

	var unpacked strings.Builder
	var prevR rune
	isPrevEscaped := false
	lastPosition := utf8.RuneCountInString(packed) - 1
	for i, r := range runes {
		switch {
		case unicode.IsDigit(r):
			if !isPrevEscaped && prevR == 0 {
				return "", ErrNumbersAreNotAllowed
			}

			if !isPrevEscaped && prevR == '\\' {
				prevR = r
				isPrevEscaped = true
				continue
			}

			digit, _ := strconv.Atoi(string(r))
			unpacked.WriteString(strings.Repeat(string(prevR), digit))
			prevR = 0
		case r == '\\':
			if !isPrevEscaped && prevR == '\\' {
				prevR = r
				isPrevEscaped = true
				continue
			} else if i == lastPosition {
				return "", ErrInvalidEscaping
			}

			unpacked.WriteRune(prevR)
			prevR = r
		default:
			if !isPrevEscaped && prevR == '\\' {
				return "", ErrInvalidEscaping
			}

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
