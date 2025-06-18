# gospark

A powerful command-line tool for generating beautiful Unicode sparklines in your terminal. Transform your numeric data into compact visual representations perfect for dashboards, monitoring, and data analysis.

![Version](https://img.shields.io/badge/version-1.0.0-blue)
![License](https://img.shields.io/badge/license-MIT-green)
![Go Version](https://img.shields.io/badge/go-1.24+-blue)

## ✨ Features

- **Multiple Input Formats**: Command-line arguments, stdin piping, or mixed separators (space, comma, pipe, semicolon)
- **Visual Modes**: Horizontal and vertical sparklines with reverse ordering
- **Rich Statistics**: Sum, min/max, and average calculations
- **Full Color Support**: Background and foreground colors with 8 standard terminal colors
- **Robust Error Handling**: Comprehensive validation for edge cases and overflow protection
- **High Performance**: Optimized for large datasets with built-in benchmarking
- **Zero Dependencies**: Single binary with no external requirements

## 📦 Installation

### Pre-built Binaries
Download the latest release from the [releases page](https://github.com/khchehab/gospark/releases).

### Build from Source
```bash
git clone https://github.com/khchehab/gospark.git
cd gospark
go build -o gospark cmd/main.go
```

### Go Install
```bash
go install github.com/khchehab/gospark/cmd@latest
```

## 🚀 Quick Start

```bash
# Basic sparkline
$ gospark 1 5 22 13 53
▁▂█▅█

# With statistics
$ gospark 1 5 22 13 53 --stats
▁▂█▅█ (min:1 max:53 avg:18.80)

# From piped data
$ echo "1 2 3 4 5" | gospark
▁▂▄▆█

# Vertical layout
$ gospark 1 2 3 4 5 --vertical
▏
▎
▌
▊
█
```

## 📖 Usage

```
gospark [flags]... value...

Flags:
  -b, --bgcolor string   background color of the sparkline graph
  -f, --fgcolor string   foreground color of the sparkline graph  
  -r, --reverse          reverse the graph
  -s, --sum              show sum of points
  -t, --stats            show stats (min, max and avg)
  -v, --vertical         show vertical graph
  -h, --help             help for gospark
      --version          version for gospark
```

## 🎨 Examples

### Basic Usage

```bash
# Simple horizontal sparkline
$ gospark 1 5 22 13 5
▁▂█▅▂

# Multiple arguments
$ gospark 0 30 55 80 33 150
▁▂▃▄▂█

# Single values use middle tick
$ gospark 42
▅
```

### Input Formats

```bash
# Space-separated values
$ gospark 1 2 3 4 5
▁▂▄▆█

# Comma-separated in quotes
$ gospark "1,2,3,4,5"
▁▂▄▆█

# Pipe-separated
$ gospark "1|2|3|4|5"
▁▂▄▆█

# Mixed separators
$ gospark "1,2 3|4;5"
▁▂▄▆█

# Stdin input
$ echo "9 13 5 17 1" | gospark
▄▆▂█▁

# From file
$ cat numbers.txt | gospark

# Command output
$ ps aux | awk '{print $3}' | gospark --stats
```

### Visual Modes

```bash
# Horizontal (default)
$ gospark 1 2 3 4 5
▁▂▄▆█

# Vertical layout
$ gospark 1 2 3 4 5 --vertical
▏
▎
▌
▊
█

# Reversed order
$ gospark 1 2 3 4 5 --reverse
█▆▄▂▁

# Combined vertical + reverse
$ gospark 1 2 3 4 5 --vertical --reverse
█
▊
▌
▎
▏
```

### Statistics and Summaries

```bash
# Show sum
$ gospark 1 2 3 4 5 --sum
▁▂▄▆█ (sum:15)

# Show detailed statistics
$ gospark 1 2 3 4 5 --stats
▁▂▄▆█ (min:1 max:5 avg:3.00)

# Combined sum and stats
$ gospark 1 2 3 4 5 --sum --stats
▁▂▄▆█ (sum:15 min:1 max:5 avg:3.00)

# With negative numbers
$ gospark -- -5 -1 0 1 5 --stats
▁▃▄▅█ (min:-5 max:5 avg:0.00)
```

### Color Support

```bash
# Red background
$ gospark 1 2 3 4 5 --bgcolor red
[colored output]

# Blue text
$ gospark 1 2 3 4 5 --fgcolor blue
[colored output]

# Combined colors
$ gospark 1 2 3 4 5 --bgcolor red --fgcolor white
[colored output]

# Available colors
# black, red, green, yellow, blue, magenta, cyan, white
```

### Advanced Examples

```bash
# System monitoring
$ while true; do
    cpu=$(top -l 1 | grep "CPU usage" | awk '{print $3}' | sed 's/%//')
    echo "$cpu" | gospark --bgcolor green --stats
    sleep 1
  done

# Network traffic visualization  
$ netstat -b | awk '{print $2}' | gospark --vertical --fgcolor cyan

# File sizes in directory
$ ls -la | awk '{print $5}' | grep -v '^$' | gospark --sum

# Git commit frequency
$ git log --format="%ad" --date=short | sort | uniq -c | awk '{print $1}' | gospark --stats

# Stock price changes
$ curl -s 'api.example.com/stocks' | jq '.prices[]' | gospark --fgcolor green
```

### Error Handling

```bash
# Invalid numbers
$ gospark abc 123
Error: invalid number: abc

# Unsupported special values
$ gospark inf
Error: infinite numbers not supported: inf

$ gospark nan  
Error: NaN (not a number) not supported: nan

# Numbers too large
$ gospark 9223372036854775808
Error: number is too large: 9223372036854775808

# Overflow protection
$ gospark 9223372036854775807 9223372036854775807 --sum
Error: numbers are too large, sum would overflow

# Invalid colors
$ gospark 1 2 3 --bgcolor purple
Error: invalid color: purple
```

## 🔧 Technical Details

### Supported Number Formats
- **Integers**: `-123`, `0`, `456`
- **Floats**: `1.5`, `-2.7`, `3.14159` (converted to integers)
- **Scientific Notation**: `1e3`, `2.5e-1` (within int64 range)
- **Range**: `-9,223,372,036,854,775,808` to `9,223,372,036,854,775,807`

### Unicode Characters
- **Horizontal**: `▁▂▃▄▅▆▇█` (8 levels)
- **Vertical**: `▏▎▍▌▋▊▉█` (8 levels)
- **Same Values**: Uses middle characters (`▅` or `▋`)

### Color Codes
Supports standard 8-color terminal palette:
- `black`, `red`, `green`, `yellow`, `blue`, `magenta`, `cyan`, `white`

### Performance
- **Memory Efficient**: Processes data in single pass
- **Overflow Protected**: Safe handling of large number sums  
- **Fast Processing**: Optimized for datasets with millions of points
- **Comprehensive Testing**: 95+ test cases covering edge cases

## 📊 Real-World Use Cases

### Use Cases
```bash
# Response time monitoring  
$ curl -w "%{time_total}\n" -s -o /dev/null example.com | gospark

# Sales data visualization
$ cat sales.csv | cut -d, -f3 | gospark --sum --fgcolor green

# Temperature readings
$ cat weather.log | grep temp | awk '{print $3}' | gospark --vertical

# Error rate trends  
$ grep ERROR app.log | wc -l | gospark --bgcolor yellow
```

### Development Setup
```bash
git clone https://github.com/khchehab/gospark.git
cd gospark
go mod download
go test ./...
go build -o gospark cmd/main.go
```

### Running Tests
```bash
# Run all tests
go test -v

# Run with race detection
go test -race

# Run benchmarks
go test -bench=.

# Test coverage
go test -cover
```

## 📄 License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## 🙏 Acknowledgments

- Inspired by [Zach Holman's spark](https://github.com/holman/spark) shell script
- Built with [Cobra CLI](https://github.com/spf13/cobra) framework
- Unicode block characters from the [Unicode Standard](https://unicode.org/charts/PDF/U2580.pdf)

---

**Made with ❤️ in Go** | **Star ⭐ if this project helped you!**
