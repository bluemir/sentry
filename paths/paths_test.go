package paths

import (
	"fmt"
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
		result := New(tt.in).Glob()

		if fmt.Sprintf("%v", result) != fmt.Sprintf("%v", tt.out) {
			t.Errorf("glob(%q) => %q, want %q", tt.in, result, tt.out)
		}
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
			[]string{"a", "a", "b", "b", "c", "c", "d", "d"},
		},
	}

	for _, tt := range tests {
		result := New([]string{"a", "b", "c", "d"}).Expand(tt.in)

		if fmt.Sprintf("%v", result) != fmt.Sprintf("%v", tt.out) {
			t.Errorf("expand(%q) => %q, want %q", tt.in, result, tt.out)
		}
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
		result := New([]string{"a", "b", "c", "d"}).Filter(tt.in)

		if fmt.Sprintf("%v", result) != fmt.Sprintf("%v", tt.out) {
			t.Errorf("filter(%q) => %q, want %q", tt.in, result, tt.out)
		}
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
		result := New([]string{"a", "b", "c", "d"}).Map(tt.in)

		if fmt.Sprintf("%v", result) != fmt.Sprintf("%v", tt.out) {
			t.Errorf("Map(%q) => %q, want %q", tt.in, result, tt.out)
		}
	}
}
