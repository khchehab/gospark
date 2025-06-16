package spark

import "fmt"

var (
	ColorMap = map[string]int{
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

func ValidateColor(color string) error {
	if color == "" {
		return nil
	}

	if _, exists := ColorMap[color]; !exists {
		return fmt.Errorf("invalid color: %s", color)
	}

	return nil
}
