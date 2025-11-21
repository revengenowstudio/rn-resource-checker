package utils

import "strings"

func SplitArgString(input string) []string {
	sep := " "
	switch {
	case strings.Contains(input, ","):
		sep = ","
	case strings.Contains(input, ";"):
		sep = ";"
	}

	return strings.Split(input, sep)
}
