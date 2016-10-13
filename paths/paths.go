package paths

import (
	log "github.com/Sirupsen/logrus"
	"path/filepath"
)

type Paths []string

func New(strs []string) Paths {
	return Paths(strs)
}

func (paths Paths) Glob() Paths {
	return paths.Expand(func(p string) []string {
		if result, err := filepath.Glob(p); err != nil {
			panic(err)
		} else {
			return result
		}
	})

}
func (paths Paths) Expand(expander func(string) []string) Paths {
	var result []string
	for _, str := range paths {
		result = append(result, expander(str)...)
	}
	log.Debug(result)
	return result
}
func (paths Paths) Filter(filter func(string) bool) Paths {
	var result []string
	for _, str := range paths {
		if filter(str) {
			result = append(result, str)
		}
	}
	log.Debug(result)
	return result
}
func (paths Paths) Map(mapper func(string) string) Paths {
	result := make([]string, len(paths))
	for k, str := range paths {
		result[k] = mapper(str)
	}
	log.Debug(result)
	return result
}
