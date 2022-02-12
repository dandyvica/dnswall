package main

import (
	"flag"
	"fmt"
	"log"
	"os"

	//"os"
	"strconv"
	"strings"
)

const Usage = `
NAME
	goget: this is a grep utility written in Go. Project repository: https://github.com/dandyvica/gogrep
	Files matching the regex (see regexp Go syntax: https://golang.org/pkg/regexp/syntax) are displayed in
	addition to matching lines, with different colors.

USAGE
	goget [OPTIONS...] PATTERN [FILE...]

OPTIONS
	-u, -url
		base url

	-b, -bounds
		bounds expressed as lower..upper

	-p, -padding
		padding value (leading 0's)
`

type Bound struct {
	lower     int
	upper     int
	maxLength int
}

// this will hold all options
type CliOptions struct {
	// url base from which download
	url string

	// bounds in url
	bound Bound

	// whether we pad the bounds with leading 0: e.g.: if bounds == 2..10, do we loop on: 02 03 ... 10 ?
	// padding = the bound length as a string. e.g.: if padding == 3 and bounds == 2..20, bounds start at 002
	padding int
}

func CliArgs() CliOptions {
	// init struct
	var options CliOptions
	var bounds string

	// if set, we want the line number from the file
	flag.StringVar(&options.url, "u", "", "base url to fetch")
	flag.StringVar(&bounds, "b", "", "bounds expressed like lower..upper (e.g.: 4..10)")
	flag.IntVar(&options.padding, "p", 0, "padding length")

	flag.Usage = func() {
		fmt.Print(Usage)
	}

	flag.Parse()

	// nothing entered: show help
	// if len(flag.Args()) == 0 {
	// 	fmt.Print(Usage)
	// 	os.Exit(0)
	// }

	// mandatory arguments
	if options.url == "" {
		log.Fatalf("argument url is mandatory")
		os.Exit(1)
	}

	if bounds == "" {
		log.Fatalf("argument bound is mandatory")
		os.Exit(2)
	}

	// some checks
	if options.padding < 0 {
		log.Fatalf("padding should be positive!")
		os.Exit(3)
	}

	// get bounds
	splittedBounds := strings.Split(bounds, "..")
	options.bound.lower = toInt(splittedBounds[0])
	options.bound.upper = toInt(splittedBounds[1])
	options.bound.maxLength = Max(len(splittedBounds[0]), len(splittedBounds[1]))

	return options
}

// convert a boolean to 0 or 1
func toInt(s string) int {
	v, err := strconv.Atoi(s)
	if err != nil {
		log.Fatalf("unable to convert %s value to integer", s)
	}
	return v
}

// Max returns the larger of x or y.
func Max(x, y int) int {
	if x < y {
		return y
	}
	return x
}
