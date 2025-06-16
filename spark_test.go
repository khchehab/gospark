package spark

import "testing"

var testCases = []struct {
	name      string
	args      []int
	bgColor   string
	fgColor   string
	showSum   bool
	showStats bool
	expected  string
}{
	// Basic functionality tests
	{"empty args", nil, "", "", false, false, ""},
	{"random args #1", []int{1, 5, 22, 13, 5}, "", "", false, false, "▁▂█▅▂"},
	{"random args #2", []int{0, 30, 55, 80, 33, 150}, "", "", false, false, "▁▂▃▄▂█"},
	{"random args #3", []int{5, 20}, "", "", false, false, "▁█"},
	{"small and very large numbers", []int{1, 2, 3, 4, 100, 5, 10, 20, 50, 300}, "", "", false, false, "▁▁▁▁▃▁▁▁▂█"},
	{"one, fifty and hundred", []int{1, 50, 100}, "", "", false, false, "▁▄█"},
	{"two, four, eight", []int{2, 4, 8}, "", "", false, false, "▁▃█"},
	{"one to five", []int{1, 2, 3, 4, 5}, "", "", false, false, "▁▂▄▆█"},
	{"same number", []int{1, 1, 1, 1}, "", "", false, false, "▅▅▅▅"},

	// Color tests
	{"one to five with blue background", []int{1, 2, 3, 4, 5}, "blue", "", false, false, "\033[44m▁▂▄▆█\033[0m"},
	{"one to five with red foreground", []int{1, 2, 3, 4, 5}, "", "red", false, false, "\033[31m▁▂▄▆█\033[0m"},
	{"one to five with blue background and red foreground", []int{1, 2, 3, 4, 5}, "blue", "red", false, false, "\033[44;31m▁▂▄▆█\033[0m"},

	// Sum tests
	{"simple numbers with sum", []int{1, 2, 3, 4, 5}, "", "", true, false, "▁▂▄▆█ (sum:15)"},
	{"zeros with sum", []int{0, 1, 2}, "", "", true, false, "▁▄█ (sum:3)"},
	{"negative numbers with sum", []int{-5, -1, 0, 1, 5}, "", "", true, false, "▁▃▄▅█ (sum:0)"},
	{"single number with sum", []int{42}, "", "", true, false, "▅ (sum:42)"},
	{"same numbers with sum", []int{5, 5, 5, 5}, "", "", true, false, "▅▅▅▅ (sum:20)"},

	// Stats tests
	{"simple numbers with stats", []int{1, 2, 3, 4, 5}, "", "", false, true, "▁▂▄▆█ (min:1 max:5 avg:3.00)"},
	{"zeros with stats", []int{0, 1, 2}, "", "", false, true, "▁▄█ (min:0 max:2 avg:1.00)"},
	{"negative numbers with stats", []int{-5, -1, 0, 1, 5}, "", "", false, true, "▁▃▄▅█ (min:-5 max:5 avg:0.00)"},
	{"single number with stats", []int{42}, "", "", false, true, "▅ (min:42 max:42 avg:42.00)"},
	{"same numbers with stats", []int{5, 5, 5, 5}, "", "", false, true, "▅▅▅▅ (min:5 max:5 avg:5.00)"},
	{"decimal average with stats", []int{1, 2, 4}, "", "", false, true, "▁▃█ (min:1 max:4 avg:2.33)"},

	// Sum and Stats combined tests
	{"simple numbers with sum and stats", []int{1, 2, 3, 4, 5}, "", "", true, true, "▁▂▄▆█ (sum:15 min:1 max:5 avg:3.00)"},
	{"negative numbers with sum and stats", []int{-2, -1, 0, 1, 2}, "", "", true, true, "▁▂▄▆█ (sum:0 min:-2 max:2 avg:0.00)"},
	{"single number with sum and stats", []int{10}, "", "", true, true, "▅ (sum:10 min:10 max:10 avg:10.00)"},

	// Colors with Sum tests
	{"blue background with sum", []int{1, 2, 3}, "blue", "", true, false, "\033[44m▁▄█\033[0m (sum:6)"},
	{"red foreground with sum", []int{1, 2, 3}, "", "red", true, false, "\033[31m▁▄█\033[0m (sum:6)"},
	{"blue bg and red fg with sum", []int{1, 2, 3}, "blue", "red", true, false, "\033[44;31m▁▄█\033[0m (sum:6)"},

	// Colors with Stats tests
	{"blue background with stats", []int{1, 2, 3}, "blue", "", false, true, "\033[44m▁▄█\033[0m (min:1 max:3 avg:2.00)"},
	{"red foreground with stats", []int{1, 2, 3}, "", "red", false, true, "\033[31m▁▄█\033[0m (min:1 max:3 avg:2.00)"},
	{"blue bg and red fg with stats", []int{1, 2, 3}, "blue", "red", false, true, "\033[44;31m▁▄█\033[0m (min:1 max:3 avg:2.00)"},

	// All options combined tests
	{"blue background with sum and stats", []int{1, 2, 3}, "blue", "", true, true, "\033[44m▁▄█\033[0m (sum:6 min:1 max:3 avg:2.00)"},
	{"red foreground with sum and stats", []int{1, 2, 3}, "", "red", true, true, "\033[31m▁▄█\033[0m (sum:6 min:1 max:3 avg:2.00)"},
	{"all options combined", []int{1, 2, 3, 4, 5}, "blue", "red", true, true, "\033[44;31m▁▂▄▆█\033[0m (sum:15 min:1 max:5 avg:3.00)"},
}

func TestSpark(t *testing.T) {
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			config := &Config{
				BgColor:   tc.bgColor,
				FgColor:   tc.fgColor,
				ShowSum:   tc.showSum,
				ShowStats: tc.showStats,
			}
			actual := Spark(tc.args, config)

			if actual != tc.expected {
				t.Errorf("got '%s', want '%s'", actual, tc.expected)
			}
		})
	}
}

func BenchmarkSparkWithoutColors(b *testing.B) {
	for i := 0; i < b.N; i++ {
		Spark([]int{1, 5, 22, 13, 5}, &Config{
			BgColor:   "",
			FgColor:   "",
			ShowSum:   false,
			ShowStats: false,
		})
	}
}

func BenchmarkSparkWithBackground(b *testing.B) {
	for i := 0; i < b.N; i++ {
		Spark([]int{1, 5, 22, 13, 5}, &Config{
			BgColor:   "red",
			FgColor:   "",
			ShowSum:   false,
			ShowStats: false,
		})
	}
}

func BenchmarkSparkWithForeground(b *testing.B) {
	for i := 0; i < b.N; i++ {
		Spark([]int{1, 5, 22, 13, 5}, &Config{
			BgColor:   "red",
			FgColor:   "blue",
			ShowSum:   false,
			ShowStats: false,
		})
	}
}

func BenchmarkSparkWithBackgroundAndForeground(b *testing.B) {
	for i := 0; i < b.N; i++ {
		Spark([]int{1, 5, 22, 13, 5}, &Config{
			BgColor:   "red",
			FgColor:   "blue",
			ShowSum:   false,
			ShowStats: false,
		})
	}
}

func BenchmarkSparkWithSum(b *testing.B) {
	for i := 0; i < b.N; i++ {
		Spark([]int{1, 5, 22, 13, 5}, &Config{
			BgColor:   "",
			FgColor:   "",
			ShowSum:   true,
			ShowStats: false,
		})
	}
}

func BenchmarkSparkWithStats(b *testing.B) {
	for i := 0; i < b.N; i++ {
		Spark([]int{1, 5, 22, 13, 5}, &Config{
			BgColor:   "",
			FgColor:   "",
			ShowSum:   false,
			ShowStats: true,
		})
	}
}

func BenchmarkSparkWithSumAndStats(b *testing.B) {
	for i := 0; i < b.N; i++ {
		Spark([]int{1, 5, 22, 13, 5}, &Config{
			BgColor:   "",
			FgColor:   "",
			ShowSum:   true,
			ShowStats: true,
		})
	}
}
