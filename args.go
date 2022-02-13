package main

import (
	"flag"
	"fmt"
	"log"
	"os"
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

// this will hold all options
type CliOptions struct {
	// DNS resolver to which forward requests
	resolver string

	// Whole resolver address (e.g.: 1.1.1.1:53)
	resolverAddress string

	// log file
	logFile string

	// do not filter, just log requests
	noFilter bool

	// configuration file
	configFile string

	// pointer on log file
	logFileHAndle *os.File

	// debug flag
	debug bool
}

func CliArgs() CliOptions {
	// init struct
	var options CliOptions

	// if set, we want the line number from the file
	flag.StringVar(&options.resolver, "r", "1.1.1.1", "DNS resolver to which unfiltered requests are forwarded")
	flag.StringVar(&options.logFile, "l", "dnswall.log", "log file name and path")
	flag.StringVar(&options.configFile, "c", "dnswall.yml", "configuration file name and path")
	flag.BoolVar(&options.noFilter, "n", false, "don't filter DNS requests")
	flag.BoolVar(&options.debug, "d", false, "debug flag")

	flag.Usage = func() {
		fmt.Print(Usage)
	}

	flag.Parse()

	// open or create log file
	f, err := os.OpenFile(options.logFile, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		log.Fatalf("error: <%v> opening file: <%v>", options.logFile, err)
	}
	log.SetOutput(f)

	// customize log format
	log.SetFlags(log.Ldate | log.Lmicroseconds)

	// save handler to close it gracefully when exiting
	options.logFileHAndle = f

	// build resolver full address
	options.resolverAddress = fmt.Sprintf("%s:53", options.resolver)

	// nothing entered: show help
	// if len(flag.Args()) == 0 {
	// 	fmt.Print(Usage)
	// 	os.Exit(0)
	// }

	return options
}
