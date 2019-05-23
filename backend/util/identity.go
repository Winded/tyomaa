package util

import (
	"regexp"
)

var (
	NameIdentifierError  = "Can only contain lower- and uppercase letters, numbers, dashes or forward slashes"
	nameIdentifierRegExp *regexp.Regexp
)

func init() {
	nameIdentifierRegExp = regexp.MustCompile(`^[a-zA-Z0-9\-\/]+$`)
}

func ValidateNameIdentifier(nameId string) bool {
	return nameIdentifierRegExp.MatchString(nameId)
}
