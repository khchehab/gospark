package main

import (
	"bufio"
	"fmt"
	"github.com/spf13/cobra"
	"os"
	"strconv"
	"strings"
)

const (
	Version = "0.1.0"
)

func main() {
	rootCmd := &cobra.Command{
		Use:   "spark VALUE...",
		Short: "Generate sparkline charts from numbers",
		Long: `Generate sparkline charts from numeric data provided as command line arguments or piped input.
Sparklines are small, word-sized graphics that show data trends without axes or coordinates.`,
		Version: Version,
		Example: `  spark 1 5 22 13 53       => ▁▁▃▂█
  spark 0,30,55,80,33,150  => ▁▂▃▄▂█
  echo 9 13 5 17 1 | spark => ▄▆▂█▁`,
		RunE: func(_ *cobra.Command, args []string) error {
			numbers, err := validateArgs(args, os.Stdin)
			if err != nil {
				return err
			}

			sparks := spark(numbers)
			fmt.Println(sparks)

			return nil
		},
	}

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
		flattened = append(flattened, strings.Fields(s)...)
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

func spark(data []int) string {
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

	var ticks []string
	if minimum == maximum {
		ticks = []string{"▅", "▆"}
	} else {
		ticks = []string{"▁", "▂", "▃", "▄", "▅", "▆", "▇", "█"}
	}

	f := ((maximum - minimum) << 8) / (len(ticks) - 1)
	if f < 1 {
		f = 1
	}

	sparklines := make([]string, len(data))
	for i, n := range data {
		sparklines[i] = ticks[((n-minimum)<<8)/f]
	}

	return strings.Join(sparklines, "")
}
