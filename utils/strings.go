package utils

import (
	log "github.com/Sirupsen/logrus"
)

type Strings []string

func NewStrings(strs ...string) Strings {
	return strs
}

func (ss Strings) Map(mapper func(string) string) Strings {
	result := make([]string, len(ss))
	for k, str := range ss {
		result[k] = mapper(str)
	}
	log.Debug(result)
	return result
}
func (ss Strings) Filter(filter func(string) bool) Strings {
	result := Strings{}
	for _, str := range ss {
		if filter(str) {
			result = append(result, str)
		}
	}
	return result
}
func (ss Strings) Expand(expander func(string) []string) Strings {
	result := Strings{}
	for _, str := range ss {
		result = append(result, expander(str)...)
	}
	return result
}
