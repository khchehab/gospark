package main

import (
	"bufio"
	"fmt"
	"github.com/spf13/cobra"
	"os"
	"strconv"
	"strings"
	"unicode"
)

const (
	Version = "0.2.0"
)

var (
	colorMap = map[string]int{
		"black":   0,
		"red":     1,
		"green":   2,
		"yellow":  3,
		"blue":    4,
		"magenta": 5,
		"cyan":    6,
		"white":   7,
	}
)

func main() {
	var bgColor string
	var fgColor string

	rootCmd := &cobra.Command{
		Use:   "spark VALUE...",
		Short: "Generate sparkline charts from numbers",
		Long: `Generate sparkline charts from numeric data provided as command line arguments or piped args.
Sparklines are small, word-sized graphics that show data trends without axes or coordinates.

Sparklines can be colored (background and foreground) with a list of predefined color names:
black, red, green, yellow, blue, magenta, cyan and white.`,
		Version: Version,
		Example: `  spark 1 5 22 13 53       => ▁▁▃▂█
  spark 0,30,55,80,33,150  => ▁▂▃▄▂█
  echo 9 13 5 17 1 | spark => ▄▆▂█▁`,
		RunE: func(_ *cobra.Command, args []string) error {
			if err := validateColor(bgColor); err != nil {
				return err
			}

			if err := validateColor(fgColor); err != nil {
				return err
			}

			numbers, err := validateArgs(args, os.Stdin)
			if err != nil {
				return err
			}

			sparks := spark(numbers, bgColor, fgColor)
			fmt.Println(sparks)

			return nil
		},
	}

	rootCmd.Flags().StringVarP(&bgColor, "bgcolor", "b", "", "background color of the sparkline graph")
	rootCmd.Flags().StringVarP(&fgColor, "fgcolor", "f", "", "foreground color of the sparkline graph")

	if rootCmd.Execute() != nil {
		os.Exit(1)
	}
}

func validateArgs(args []string, file *os.File) ([]int, error) {
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

	numbers := make([]int, 0, len(flattened))
	for _, n := range flattened {
		f, err := strconv.ParseFloat(n, 64)
		if err != nil {
			return nil, fmt.Errorf("%s is not a number", n)
		}

		numbers = append(numbers, int(f))
	}

	return numbers, nil
}

func spark(data []int, bgColor, fgColor string) string {
	if len(data) == 0 {
		return ""
	}

	minimum, maximum := data[0], data[0]

	for i := 1; i < len(data); i++ {
		if data[i] < minimum {
			minimum = data[i]
		}
		if data[i] > maximum {
			maximum = data[i]
		}
	}

	var ticks []rune
	if minimum == maximum {
		ticks = []rune{'▅', '▆'}
	} else {
		ticks = []rune{'▁', '▂', '▃', '▄', '▅', '▆', '▇', '█'}
	}

	f := ((maximum - minimum) << 8) / (len(ticks) - 1)
	if f < 1 {
		f = 1
	}

	sparklines := make([]rune, len(data))
	for i, n := range data {
		sparklines[i] = ticks[((n-minimum)<<8)/f]
	}

	prefix, suffix := "", ""
	if bgColor != "" || fgColor != "" {
		prefix = "\033["
		if bgColor != "" {
			prefix += strconv.Itoa(40 + colorMap[bgColor])
		}

		if fgColor != "" {
			if bgColor != "" {
				prefix += ";"
			}
			prefix += strconv.Itoa(30 + colorMap[fgColor])
		}
		prefix += "m"

		suffix = "\033[0m"
	}

	return prefix + string(sparklines) + suffix
}

func validateColor(color string) error {
	if color == "" {
		return nil
	}

	if _, exists := colorMap[color]; !exists {
		return fmt.Errorf("invalid color: %s", color)
	}

	return nil
}

func isSeparator(r rune) bool {
	return unicode.IsSpace(r) || r == ',' || r == '|'
}
