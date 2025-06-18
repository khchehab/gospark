package main

import (
	"fmt"
	"github.com/spf13/cobra"
	spark "gospark"
	"os"
)

const (
	Version = "1.0.0"
)

func main() {
	config := &spark.Config{}

	rootCmd := &cobra.Command{
		Use:                   "spark [flags]... value...",
		DisableFlagsInUseLine: true,
		Short:                 "Draw sparkline charts from numbers",
		Long: `Draw sparkline charts from numeric data provided as command line arguments or piped data.
Sparklines are small, word-sized graphics that show data trends without axes or coordinates.

Numbers can be separated by any space character, comma, pipe (|) or semi-colon.
For negative numbers, use the flag separator '--' (flags must come before it).

Sparklines can be colored (background and foreground) with a list of predefined color names:
black, red, green, yellow, blue, magenta, cyan and white.`,
		Version: Version,
		Example: `  spark 1 5 22 13 53               => ▁▁▃▂█
 spark 0,30,55,80,33,150 --sum    => ▁▂▃▄▂█ (sum:348)
 echo "9 13 5 17 1" | spark       => ▄▆▂█▁
 spark "1|2|3|4|5" --stats        => ▁▂▄▆█ (min:1 max:5 avg:3.00)
 spark --sum -- -5 -1 0 1 5       => ▁▃▄▅█ (sum:0)`,
		RunE: func(_ *cobra.Command, args []string) error {
			if err := config.Validate(); err != nil {
				return err
			}

			data, err := spark.ValidateArgs(args, os.Stdin)
			if err != nil {
				return err
			}

			sparks, err := spark.Spark(data, config)
			if err != nil {
				return err
			}
			fmt.Println(sparks)

			return nil
		},
	}

	rootCmd.Flags().StringVarP(&config.BgColor, "bgcolor", "b", "", "background color of the sparkline graph")
	rootCmd.Flags().StringVarP(&config.FgColor, "fgcolor", "f", "", "foreground color of the sparkline graph")
	rootCmd.Flags().BoolVarP(&config.ShowSum, "sum", "s", false, "show sum of points")
	rootCmd.Flags().BoolVarP(&config.ShowStats, "stats", "t", false, "show stats (min, max and avg)")
	rootCmd.Flags().BoolVarP(&config.Reverse, "reverse", "r", false, "reverse the graph")
	rootCmd.Flags().BoolVarP(&config.Vertical, "vertical", "v", false, "show vertical graph")

	if rootCmd.Execute() != nil {
		os.Exit(1)
	}
}
