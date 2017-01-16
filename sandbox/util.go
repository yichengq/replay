package main

import "strings"

func polishReasonBinsOutput(output []byte) string {
	return excludeLine(string(output), []string{"-impl prog.re", "Command exited with code"})
}

func excludeLine(str string, excludedWords []string) string {
	ss := strings.Split(str, "\n")
	var newss []string
	for _, s := range ss {
		var has bool
		for _, w := range excludedWords {
			if strings.Contains(s, w) {
				has = true
			}
		}
		if !has {
			newss = append(newss, s)
		}
	}
	return strings.Join(newss, "\n")
}
