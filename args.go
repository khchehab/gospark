package spark

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
	"unicode"
)

func ValidateArgs(args []string, stdin *os.File) ([]int, error) {
	hasArgs := len(args) > 0

	stdinStat, _ := stdin.Stat()
	hasStdinData := (stdinStat.Mode() & os.ModeCharDevice) == 0

	var source []string
	if hasArgs {
		source = args
	} else if hasStdinData {
		scanner := bufio.NewScanner(stdin)
		for scanner.Scan() {
			source = append(source, scanner.Text())
		}
	}

	return parseSource(source)
}

func parseSource(source []string) ([]int, error) {
	var flattened []string
	for _, s := range source {
		flattened = append(flattened, strings.FieldsFunc(s, isSeparator)...)
	}

	if len(flattened) < 1 {
		return nil, fmt.Errorf("no numeric data provided - specify numbers as arguments or pipe data via stdin")
	}

	data := make([]int, 0, len(flattened))
	for _, n := range flattened {
		f, err := strconv.ParseFloat(n, 64)
		if err != nil {
			return nil, fmt.Errorf("invalid number: %s", n)
		}

		data = append(data, int(f))
	}

	return data, nil
}

func isSeparator(r rune) bool {
	return unicode.IsSpace(r) || r == ',' || r == '|' || r == ';'
}
