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

	prevStr, symbols := rune(0), make([]string, 0, len(str))

	for i, r := range str {
		if unicode.IsNumber(r) && i == 0 || unicode.IsDigit(prevStr) && unicode.IsDigit(r) {
			return "", ErrInvalidString
		}

		prevStr = r

		if !unicode.IsNumber(r) {
			symbols = append(symbols, string(r))

			continue
		}

		previousSymbol := symbols[len(symbols)-1]
		intVal, _ := strconv.Atoi(string(r))

		symbols[len(symbols)-1] = strings.Repeat(previousSymbol, intVal)
	}

	return strings.Join(symbols, ""), nil
}
