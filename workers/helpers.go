package workers

import "strings"

func StringMatch(msg string, subs []string) (match bool){
	for _, sub := range subs {
		if strings.Contains(strings.ToLower(msg), strings.ToLower(sub)) {
			match = true
			break
		}
	}
	return match
}

func RegexMatch(msg string, regex []string) (match bool){
	for _, sub := range regex {
		if strings.Contains(strings.ToLower(msg), strings.ToLower(sub)) {
			match = true
			break
		}
	}
	return match
}