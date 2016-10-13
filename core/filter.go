package core

import (
	"path/filepath"

	log "github.com/Sirupsen/logrus"
)

type fileNameFilter struct {
	patterns []string
}

func newFileNameFilter(pattern []string) *fileNameFilter {

	return &fileNameFilter{
		patterns: pattern,
	}
}

func (filter *fileNameFilter) check(path string) bool {

	for _, pattern := range filter.patterns {
		if result, err := filepath.Match(pattern, path); err == nil && result {
			return false
		} else if err != nil {
			log.Warnf("Bad pattern : %s", pattern)
		}
	}
	return true
}
