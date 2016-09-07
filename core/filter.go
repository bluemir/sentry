package core

import (
	"regexp"

	log "github.com/Sirupsen/logrus"
)

type fileNameFilter struct {
	re *regexp.Regexp
}

func newFileNameFilter(pattern string) *fileNameFilter {
	re, err := regexp.Compile(pattern)
	if err != nil {
		log.Warn("Fail to compile filter pattern")
	}
	return &fileNameFilter{
		re,
	}
}

func (filter *fileNameFilter) check(path string) bool {
	if filter.re == nil {
		return true
	}
	return filter.re.MatchString(path)
}
