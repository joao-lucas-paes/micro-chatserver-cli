package rules

import (
	"regexp"
	"strconv"
)

func loginMatch(login string) (string, int, bool) {
	nick := ""
	channel := -1
	
	re := regexp.MustCompile(`^(\w+)\s+(\d+)$`)
	m := re.FindStringSubmatch(login)

	if re.MatchString(login) {
		nick = m[1]
		channel_val, err := strconv.Atoi(m[2])
		channel = channel_val
		if err != nil {
			return "", -1, false
		}
	}

	return nick, channel, true
}