package hw02unpackstring

import (
	"errors"
	"strconv"
	"strings"
	"unicode"
)

var ErrInvalidString = errors.New("invalid string")

func Unpack(str string) (string, error) {
	if _, err := strconv.Atoi(str); err == nil {
		return "", ErrInvalidString
	}

	var (
		b       strings.Builder
		prevStr rune
		sl      = []rune(str)
	)

	for i, r := range str {
		if unicode.IsNumber(r) && i == 0 || unicode.IsDigit(prevStr) && unicode.IsDigit(r) {
			return "", ErrInvalidString
		}

		if unicode.IsLower(r) && unicode.IsLower(prevStr) {
			b.WriteString(string(prevStr))
		} else if unicode.IsLower(prevStr) && unicode.IsDigit(r) {
			b.WriteString(strings.Repeat(string(prevStr), int(r-'0')))
		}

		if unicode.IsLower(r) && i == len(sl)-1 {
			b.WriteString(string(r))
		}

		prevStr = r
	}

	return b.String(), nil
}
