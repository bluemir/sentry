package paths

import (
	log "github.com/Sirupsen/logrus"
	"path/filepath"
)

type Paths map[string]bool

func New(strs ...string) Paths {

	result := Paths{}
	for _, str := range strs {
		result[str] = true
	}
	return result
}

func (paths Paths) Glob() Paths {
	return paths.Expand(func(p string) []string {
		if result, err := filepath.Glob(p); err != nil {
			log.Warn(err)
			return []string{p}
		} else {
			return result
		}
	})

}
func (paths Paths) Expand(expander func(string) []string) Paths {
	result := Paths{}
	for str, _ := range paths {
		e := expander(str)
		for _, s := range e {
			result[s] = true
		}
	}
	log.Debug(result)
	return result
}
func (paths Paths) Filter(filter func(string) bool) Paths {
	result := Paths{}
	for str, _ := range paths {
		if filter(str) {
			result[str] = true
		}
	}
	log.Debug(result)
	return result
}
func (paths Paths) Map(mapper func(string) string) Paths {
	result := Paths{}
	for str, _ := range paths {
		result[mapper(str)] = true
	}
	log.Debug(result)
	return result
}

func (paths Paths) Value() []string {
	var result []string

	for str, _ := range paths {
		result = append(result, str)
	}

	return result
}
