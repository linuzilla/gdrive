package utils

import (
	"regexp"
)

var QuotedStringPattern = regexp.MustCompile(`^"(.*)"`)

func StripQuotedString(source string) string {
	if matched := QuotedStringPattern.FindStringSubmatch(source); matched != nil {
		return matched[1]
	} else {
		return source
	}
}
