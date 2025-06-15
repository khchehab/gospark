package main

import "testing"

var testCases = []struct {
	name     string
	input    []int
	expected string
}{
	{"empty input", nil, ""},
	{"random input #1", []int{1, 5, 22, 13, 5}, "▁▂█▅▂"},
	{"random input #2", []int{0, 30, 55, 80, 33, 150}, "▁▂▃▄▂█"},
	{"random input #3", []int{5, 20}, "▁█"},
	{"small and very large numbers", []int{1, 2, 3, 4, 100, 5, 10, 20, 50, 300}, "▁▁▁▁▃▁▁▁▂█"},
	{"one, fifty and hundred", []int{1, 50, 100}, "▁▄█"},
	{"two, four, eight", []int{2, 4, 8}, "▁▃█"},
	{"one to five", []int{1, 2, 3, 4, 5}, "▁▂▄▆█"},
	{"same number", []int{1, 1, 1, 1}, "▅▅▅▅"},
}

func TestSpark(t *testing.T) {
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			actual := spark(tc.input)

			if actual != tc.expected {
				t.Errorf("got '%s', want '%s'", actual, tc.expected)
			}
		})
	}
}

func BenchmarkSpark(b *testing.B) {
	for i := 0; i < b.N; i++ {
		spark([]int{1, 5, 22, 13, 5})
	}
}
