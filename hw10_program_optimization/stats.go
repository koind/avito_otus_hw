package hw10programoptimization

import (
	"bufio"
	"io"
	"strings"
)

//easyjson:json
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
	result := make(DomainStat)
	scanner := bufio.NewScanner(r)
	user := User{}

	for scanner.Scan() {
		line := scanner.Bytes()

		if err := user.UnmarshalJSON(line); err != nil {
			return result, err
		}

		if strings.HasSuffix(user.Email, "."+domain) {
			result[strings.ToLower(strings.SplitN(user.Email, "@", 2)[1])]++
		}
	}

	return result, nil
}
