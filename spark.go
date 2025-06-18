package spark

import (
	"fmt"
	"math"
	"slices"
	"strconv"
	"strings"
)

func Spark(data []int, config *Config) (string, error) {
	if len(data) == 0 {
		return "", nil
	}

	minimum, maximum, sum, average, err := getStats(data)
	if err != nil {
		return "", err
	}

	ticks, separator := getTicks(minimum == maximum, config)

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

	if config.Reverse {
		slices.Reverse(sparklines)
	}

	return concatenateParts(sparklines, minimum, maximum, sum, average, separator, config), nil
}

func getStats(data []int) (int, int, int, float64, error) {
	minimum, maximum := data[0], data[0]
	sum := data[0]

	for i := 1; i < len(data); i++ {
		if data[i] < minimum {
			minimum = data[i]
		}
		if data[i] > maximum {
			maximum = data[i]
		}
		// Check for overflow before adding
		if sum > 0 && data[i] > 0 && sum > math.MaxInt-data[i] {
			return 0, 0, 0, 0, fmt.Errorf("numbers are too large, sum would overflow")
		}
		if sum < 0 && data[i] < 0 && sum < math.MinInt-data[i] {
			return 0, 0, 0, 0, fmt.Errorf("numbers are too large, sum would underflow")
		}
		sum += data[i]
	}

	average := float64(sum) / float64(len(data))

	return minimum, maximum, sum, average, nil
}

func getTicks(sameMinMax bool, config *Config) ([]rune, string) {
	var pool [8]rune
	var separator string
	if config.Vertical {
		pool = [8]rune{'▏', '▎', '▍', '▌', '▋', '▊', '▉', '█'}
		separator = "\n"
	} else {
		pool = [8]rune{'▁', '▂', '▃', '▄', '▅', '▆', '▇', '█'}
		separator = ""
	}

	var ticks []rune
	if sameMinMax {
		ticks = pool[4:6]
	} else {
		ticks = pool[:]
	}
	return ticks, separator
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

func concatenateParts(sparklines []rune, minimum, maximum, sum int, average float64, separator string, config *Config) string {
	var parts []string

	prefix, suffix := getPrefixAndSuffix(config)
	finalSparklines := make([]string, len(sparklines))
	for i, r := range sparklines {
		finalSparklines[i] = fmt.Sprintf("%s%c%s", prefix, r, suffix)
	}
	parts = append(parts, strings.Join(finalSparklines, separator))

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
