package validate

import (
	"regexp"
	"unicode"
)

var emailRegex = regexp.MustCompile("^[a-zA-Z0-9.!#$%&'*+/=?^_`{|}~-]+@[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?(?:\\.[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?)*$")

func IsEMail(email string) bool {
	return emailRegex.MatchString(email)
}

func IsDigits(num string) bool {
	for _, c := range num {
		if !unicode.IsDigit(c) {
			return false
		}
	}

	return true
}
