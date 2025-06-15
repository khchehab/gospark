package main

import "testing"

var testCases = []struct {
	name     string
	args     []int
	bgColor  string
	fgColor  string
	expected string
}{
	{"empty args", nil, "", "", ""},
	{"random args #1", []int{1, 5, 22, 13, 5}, "", "", "▁▂█▅▂"},
	{"random args #2", []int{0, 30, 55, 80, 33, 150}, "", "", "▁▂▃▄▂█"},
	{"random args #3", []int{5, 20}, "", "", "▁█"},
	{"small and very large numbers", []int{1, 2, 3, 4, 100, 5, 10, 20, 50, 300}, "", "", "▁▁▁▁▃▁▁▁▂█"},
	{"one, fifty and hundred", []int{1, 50, 100}, "", "", "▁▄█"},
	{"two, four, eight", []int{2, 4, 8}, "", "", "▁▃█"},
	{"one to five", []int{1, 2, 3, 4, 5}, "", "", "▁▂▄▆█"},
	{"same number", []int{1, 1, 1, 1}, "", "", "▅▅▅▅"},

	{"one to five with blue background", []int{1, 2, 3, 4, 5}, "blue", "", "\033[44m▁▂▄▆█\033[0m"},
	{"one to five with red foreground", []int{1, 2, 3, 4, 5}, "", "red", "\033[31m▁▂▄▆█\033[0m"},
	{"one to five with blue background and red foreground", []int{1, 2, 3, 4, 5}, "blue", "red", "\033[44;31m▁▂▄▆█\033[0m"},
}

func TestSpark(t *testing.T) {
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			actual := spark(tc.args, tc.bgColor, tc.fgColor)

			if actual != tc.expected {
				t.Errorf("got '%s', want '%s'", actual, tc.expected)
			}
		})
	}
}

func BenchmarkSparkWithoutColors(b *testing.B) {
	for i := 0; i < b.N; i++ {
		spark([]int{1, 5, 22, 13, 5}, "", "")
	}
}

func BenchmarkSparkWithBackground(b *testing.B) {
	for i := 0; i < b.N; i++ {
		spark([]int{1, 5, 22, 13, 5}, "red", "")
	}
}

func BenchmarkSparkWithForeground(b *testing.B) {
	for i := 0; i < b.N; i++ {
		spark([]int{1, 5, 22, 13, 5}, "", "blue")
	}
}

func BenchmarkSparkWithBackgroundAndForeground(b *testing.B) {
	for i := 0; i < b.N; i++ {
		spark([]int{1, 5, 22, 13, 5}, "red", "blue")
	}
}
