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
	var lastRune rune
	for _, v := range packed {
		if unicode.IsDigit(v) {
			if unicode.IsDigit(lastRune) {
				return "", fmt.Errorf("%w: number are not accepted", ErrInvalidString)
			}

			if lastRune == 0 {
				return "", fmt.Errorf("%w: digit can not be before a substring", ErrInvalidString)
			}
			digit, _ := strconv.Atoi(string(v))
			unpacked.WriteString(strings.Repeat(string(lastRune), digit))
			lastRune = 0
		} else {
			if lastRune != 0 {
				unpacked.WriteRune(lastRune)
			}
			lastRune = v
		}
	}
	if lastRune != 0 {
		unpacked.WriteRune(lastRune)
	}
	return unpacked.String(), nil
}
