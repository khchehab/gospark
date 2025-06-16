package main

import (
	"fmt"
	"github.com/spf13/cobra"
	spark "gospark"
	"os"
)

const (
	Version = "0.2.0"
)

func main() {
	config := &spark.Config{}

	rootCmd := &cobra.Command{
		Use:                   "spark [flags]... value...",
		DisableFlagsInUseLine: true,
		Short:                 "Generate sparkline charts from numbers",
		Long: `Generate sparkline charts from numeric data provided as command line arguments or piped args.
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

			sparks := spark.Spark(data, config)
			fmt.Println(sparks)

			return nil
		},
	}

	rootCmd.Flags().StringVarP(&config.BgColor, "bgcolor", "b", "", "background color of the sparkline graph")
	rootCmd.Flags().StringVarP(&config.FgColor, "fgcolor", "f", "", "foreground color of the sparkline graph")
	rootCmd.Flags().BoolVarP(&config.ShowSum, "sum", "s", false, "show sum of points")
	rootCmd.Flags().BoolVarP(&config.ShowStats, "stats", "t", false, "show stats (min, max and avg)")

	// TODO add the below options:
	// * scale and width?? how to do it with the runes used?
	// * normalize: would it change the display??
	// * smooth: how to do that if using runes??
	// * file: get numbers from a file (instead of args or pipe).
	// * output: draw the output to a file
	// * format: unicode, ascii, json or csv - how would that work?? if the end result is to draw a sparkline graph
	// * reverse: draw graph in reverse, may not be useful but can be done
	// * vertical: interesting option, need to find the appropriate runes but can be done (also should be able to handle sum and stat)
	// * help-colors: display list of colors that can be used - very interesting and useful
	// * benchmark: what to display and how would that help??
	// * raw: what is this???
	// * no-color: what is the difference with not using the fgcolor and bgcolor options??
	// * see what possible separators that can be used: add ; and see what others.

	//Great final version! Here are some beneficial options you could consider adding:
	//
	//Display & Formatting Options
	//
	//--width, -w        # Set fixed width (truncate/pad to N characters)
	//--scale, -s        # Set custom min/max range instead of auto-scaling
	//--no-color         # Disable colors even if terminal supports them
	//--raw              # Output raw Unicode without ANSI formatting
	//
	//Data Processing Options
	//
	//--normalize        # Normalize to 0-100 range before plotting
	//--smooth           # Apply simple smoothing/moving average
	//
	//Input/Output Options
	//
	//--separator, -d    # Custom separator (default: space,comma,pipe)
	//--file, -i         # Read from file instead of stdin/args
	//--output, -o       # Write to file instead of stdout
	//--format           # Output format: unicode, ascii, json, csv
	//
	//Utility Options
	//
	//--reverse          # Reverse the sparkline direction
	//--vertical         # Vertical sparkline instead of horizontal
	//--help-colors      # List all available colors
	//--benchmark        # Show performance timing info
	//
	//Most Valuable Additions (in order):
	//
	//1. --width - Very useful for consistent formatting in scripts
	//2. --stats - Adds analytical value without cluttering
	//3. --scale - Allows comparison between different datasets
	//4. --format ascii - For terminals without Unicode support
	//5. --separator - Custom delimiter support

	if rootCmd.Execute() != nil {
		os.Exit(1)
	}
}
