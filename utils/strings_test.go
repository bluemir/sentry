package utils

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

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
			[]string{"a", "a", "b", "b", "c", "c", "d", "d"},
		},
	}

	for _, tt := range tests {
		result := NewStrings("a", "b", "c", "d").Expand(tt.in)

		assert.Equal(t, result, NewStrings(tt.out...))
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
		result := NewStrings("a", "b", "c", "d").Filter(tt.in)

		assert.Equal(t, result, NewStrings(tt.out...))
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
		result := NewStrings("a", "b", "c", "d").Map(tt.in)

		assert.Equal(t, result, NewStrings(tt.out...))
	}
}
