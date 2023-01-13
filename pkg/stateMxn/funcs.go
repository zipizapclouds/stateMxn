package stateMxn

import (
	"regexp"
	"strings"
)

// https://github.com/openconfig/goyang/blob/v1.2.0/pkg/indent/indent.go#L25
// Returns s with each line in s prefixed by indent.
func identLinesInString(indent, s string) string {
	if indent == "" || s == "" {
		return s
	}
	lines := strings.SplitAfter(s, "\n")
	if len(lines[len(lines)-1]) == 0 {
		lines = lines[:len(lines)-1]
	}
	return strings.Join(append([]string{""}, lines...), indent)
}

// split str into lines of nchar characters of length, and then indent each line
func splitAndIdentLines(str string, ncharts int, indent string) string {
	var lines []string
	for len(str) > ncharts {
		lines = append(lines, str[:ncharts])
		str = str[ncharts:]
	}
	lines = append(lines, str)
	return identLinesInString(indent, strings.Join(lines, "\n"))
}

func replace2alphanum(s string) string {
	return regexp.MustCompile(`[^a-zA-Z0-9]+`).ReplaceAllString(s, "_")
}
