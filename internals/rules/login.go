package rules

import (
	"regexp"
)

func LoginMatch(login string) (string, string, bool) {
	re := regexp.MustCompile(`^(\w+)\s+(\w+)$`)
	m := re.FindStringSubmatch(login)

	if re.MatchString(login) {
		return m[1], m[2], true
	}
	return "", "", false
}