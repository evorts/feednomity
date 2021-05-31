package validate

import (
	"regexp"
	"strings"
)

var (
	usernamePattern = regexp.MustCompile("[a-zA-Z0-9.-]+")
	emailPattern = regexp.MustCompile("\\w+[._]?\\w+@\\w+.[a-zA-Z]{2,3}")
	phonePattern = regexp.MustCompile("(\\+\\d{2}|0)([1-9]+\\d{5,10})")
	passwordPattern = regexp.MustCompile("[\\w\\d@$%&^!~()]+")
	hashPattern = regexp.MustCompile("[\\w\\d@$%&^!~()]{10,128}")
)

func ValidUsername(value string) bool {
	return usernamePattern.MatchString(value)
}

func ValidEmail(value string) bool {
	return emailPattern.MatchString(value)
}

func ValidPhone(value string) bool {
	return phonePattern.MatchString(value)
}

func ValidPassword(value string) bool {
	return passwordPattern.MatchString(value)
}

func ValidHash(value string) bool {
	return hashPattern.MatchString(value)
}

func IsEmpty(value string) bool {
	if len(strings.Trim(value, " ")) < 1 {
		return true
	}
	return false
}