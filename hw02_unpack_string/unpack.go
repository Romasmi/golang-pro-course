package hw02unpackstring

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
	"unicode"
	"unicode/utf8"
)

var (
	ErrInvalidString             = errors.New("invalid string")
	NumbersAreNotAllowed         = errors.New("numbers > 9 not allowed")
	StringCanNotStartWithANumber = errors.New("string can not start with a number")
	InvalidEscaping              = errors.New("invalid escaping")
)

func Unpack(packed string) (string, error) {
	if packed == "" {
		return packed, nil
	}
	if unicode.IsDigit(rune(packed[0])) {
		return "", fmt.Errorf("%w: %w", ErrInvalidString, StringCanNotStartWithANumber)
	}

	var unpacked strings.Builder
	var prevR rune
	isPrevEscaped := false
	lastPosition := utf8.RuneCountInString(packed) - 1
	for i, r := range packed {
		switch {
		case unicode.IsDigit(r):
			if !isPrevEscaped && prevR == 0 {
				return "", fmt.Errorf("%w: %w", ErrInvalidString, NumbersAreNotAllowed)
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
			} else if prevR == '\\' && i == lastPosition {
				return "", fmt.Errorf("%w: %w", ErrInvalidString, InvalidEscaping)
			}

			unpacked.WriteRune(prevR)
			prevR = r
		default:
			if !isPrevEscaped && prevR == '\\' {
				return "", fmt.Errorf("%w: %w", ErrInvalidString, InvalidEscaping)
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
