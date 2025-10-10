package rules

import (
	"strings"
)

func LoginMatch(login string) (string, string, bool) {
	fields := strings.Fields(login) // já faz Trim e split por whitespace
	if len(fields) != 2 {
		return "", "", false
	}
	return fields[0], fields[1], true
}