package paths

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestGlob(t *testing.T) {
	var globtest = []struct {
		in  []string
		out []string
	}{
		{ //blank case
			[]string{},
			[]string{},
		},
		{
			[]string{"*.go"},
			[]string{"paths.go", "paths_test.go"},
		},
		{
			[]string{"../.gi*re"},
			[]string{"../.gitignore"},
		},
	}

	for _, tt := range globtest {
		result := New(tt.in...).Glob()

		assert.Equal(t, result, New(tt.out...))
	}
}
func TestExpand(t *testing.T) {
	var tests = []struct {
		in  func(string) []string
		out []string
	}{
		{
			func(str string) []string {
				return []string{str}
			},
			[]string{"a", "b", "c", "d"},
		},
		{
			func(str string) []string {
				return []string{str, str}
			},
			[]string{"a", "b", "c", "d"},
		},
	}

	for _, tt := range tests {
		result := New([]string{"a", "b", "c", "d"}...).Expand(tt.in)

		assert.Equal(t, result, New(tt.out...))
	}
}
func TestFilter(t *testing.T) {
	var tests = []struct {
		in  func(string) bool
		out []string
	}{
		{
			func(str string) bool {
				return true
			},
			[]string{"a", "b", "c", "d"},
		},
		{
			func(str string) bool {
				return str == "a"
			},
			[]string{"a"},
		},
	}

	for _, tt := range tests {
		result := New([]string{"a", "b", "c", "d"}...).Filter(tt.in)

		assert.Equal(t, result, New(tt.out...))
	}
}
func TestMap(t *testing.T) {
	var tests = []struct {
		in  func(string) string
		out []string
	}{
		{
			func(str string) string {
				return str
			},
			[]string{"a", "b", "c", "d"},
		},
		{
			func(str string) string {
				return str + str
			},
			[]string{"aa", "bb", "cc", "dd"},
		},
	}

	for _, tt := range tests {
		result := New([]string{"a", "b", "c", "d"}...).Map(tt.in)

		assert.Equal(t, result, New(tt.out...))
	}
}

func TestValue(t *testing.T) {
	var tests = []struct {
		in  []string
		out []string
	}{
		{
			[]string{"a", "b", "c", "d"},
			[]string{"a", "b", "c", "d"},
		},
	}
	for _, tt := range tests {
		result := New(tt.in...).Value()

		for _, str := range tt.out {
			assert.Contains(t, result, str)
		}
	}
}
