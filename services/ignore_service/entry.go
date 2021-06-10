package ignore_service

import "regexp"

type entry struct {
	Index    int
	Pattern  *regexp.Regexp
	Exactly  bool
	Positive bool
	Filename string
}
