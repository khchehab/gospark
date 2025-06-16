package spark

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
	"unicode"
)

func Spark(data []int, config *Config) string {
	if len(data) == 0 {
		return ""
	}

	minimum, maximum, sum, average := getStats(data)

	ticks, _ := getTicks(minimum == maximum, config)

	divisor := float64(maximum - minimum)
	factor := len(ticks) - 1

	sparklines := make([]rune, len(data))
	for i, n := range data {
		if divisor == 0 {
			sparklines[i] = ticks[0]
		} else {
			sparklines[i] = ticks[int(float64((n-minimum)*factor)/divisor)]
		}
	}

	return concatenateParts(sparklines, minimum, maximum, sum, average, config)
}

func ValidateArgs(args []string, file *os.File) ([]int, error) {
	hasArgs := len(args) > 0

	stat, _ := file.Stat()
	hasFileData := (stat.Mode() & os.ModeCharDevice) == 0

	var source []string
	if hasArgs {
		source = args
	} else if hasFileData {
		scanner := bufio.NewScanner(file)
		for scanner.Scan() {
			source = append(source, scanner.Text())
		}
	}

	var flattened []string
	for _, s := range source {
		flattened = append(flattened, strings.FieldsFunc(s, isSeparator)...)
	}

	if len(flattened) < 1 {
		return nil, fmt.Errorf("requires at least 1 argument")
	}

	data := make([]int, 0, len(flattened))
	for _, n := range flattened {
		f, err := strconv.ParseFloat(n, 64)
		if err != nil {
			return nil, fmt.Errorf("%s is not a number", n)
		}

		data = append(data, int(f))
	}

	return data, nil
}

func isSeparator(r rune) bool {
	return unicode.IsSpace(r) || r == ',' || r == '|' || r == ';'
}

func getStats(data []int) (int, int, int, float64) {
	minimum, maximum := data[0], data[0]
	sum := data[0]

	for i := 1; i < len(data); i++ {
		if data[i] < minimum {
			minimum = data[i]
		}
		if data[i] > maximum {
			maximum = data[i]
		}
		sum += data[i]
	}

	average := float64(sum) / float64(len(data))

	return minimum, maximum, sum, average
}

func getTicks(sameMinMax bool, config *Config) ([]rune, string) {
	var ticks []rune
	if sameMinMax {
		ticks = []rune{'▅', '▆'}
	} else {
		ticks = []rune{'▁', '▂', '▃', '▄', '▅', '▆', '▇', '█'}
	}
	return ticks, ""
}

func getPrefixAndSuffix(config *Config) (string, string) {
	if config.BgColor == "" && config.FgColor == "" {
		return "", ""
	}

	prefix := "\033["
	if config.BgColor != "" {
		prefix += strconv.Itoa(40 + ColorMap[config.BgColor])
	}

	if config.FgColor != "" {
		if config.BgColor != "" {
			prefix += ";"
		}
		prefix += strconv.Itoa(30 + ColorMap[config.FgColor])
	}
	prefix += "m"

	suffix := "\033[0m"

	return prefix, suffix
}

func concatenateParts(sparklines []rune, minimum, maximum, sum int, average float64, config *Config) string {
	var parts []string

	prefix, suffix := getPrefixAndSuffix(config)
	parts = append(parts, prefix)
	parts = append(parts, string(sparklines))
	parts = append(parts, suffix)

	if config.ShowSum || config.ShowStats {
		parts = append(parts, " (")

		var subParts []string
		if config.ShowSum {
			subParts = append(subParts, fmt.Sprintf("sum:%d", sum))
		}

		if config.ShowStats {
			subParts = append(subParts, fmt.Sprintf("min:%d", minimum))
			subParts = append(subParts, fmt.Sprintf("max:%d", maximum))
			subParts = append(subParts, fmt.Sprintf("avg:%.2f", average))
		}

		parts = append(parts, strings.Join(subParts, " "))
		parts = append(parts, ")")
	}

	return strings.Join(parts, "")
}
