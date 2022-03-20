package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"sync"

	"gopkg.in/yaml.v3"
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

// This will hold all options given from the command line
type Config struct {
	resolver        string          // DNS resolver to which forward requests
	resolverAddress string          // Whole resolver address (e.g.: 1.1.1.1:53) for the one
	timeout         int             // timeout when sending queries to resolver or sending back data to client
	logFile         string          // log file
	dontFilter      bool            // do not filter, just log requests
	yamlConfigFile  string          // configuration file
	logFileHAndle   *os.File        // pointer on log file
	debug           bool            // debug flag
	filters         FilteredDomains // list of either whitelisted domains for which DNS domain will not be blocked and blacklisted ones for which a NXDOMAIN will be sent back
	mu              sync.Mutex      // used to synchronize access to block lists
}

// This will match the YAML configuration file where all settings are defined
type YAMLConfig struct {
	Resolvers []string `yaml:"resolvers"`
	Timeout   int      `yaml:"update_timeout"`
	Filters   struct {
		Whitelist []string `yaml:"whitelist"`
		Blacklist []string `yaml:"blacklist"`
	} `yaml:"filters"`
}

// Read command line arguments and read the YAML configuration file
func readCliArgs() Config {
	// init struct
	var conf Config

	// if set, we want the line number from the file
	flag.StringVar(&conf.resolver, "r", "1.1.1.1", "DNS resolver to which unfiltered requests are forwarded")
	flag.StringVar(&conf.logFile, "l", "dnswall.log", "log file name and path")
	flag.StringVar(&conf.yamlConfigFile, "c", "dnswall.yml", "configuration file name and path")
	flag.BoolVar(&conf.dontFilter, "n", false, "don't filter DNS requests")
	flag.BoolVar(&conf.debug, "d", false, "debug flag")
	flag.IntVar(&conf.timeout, "t", 300, "timeout (in seconds) when sending queries to resolver or sending back data to client")

	flag.Usage = func() {
		fmt.Print(Usage)
	}

	flag.Parse()

	// open or create log file
	f, err := os.OpenFile(conf.logFile, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		log.Fatalf("error: <%v> opening file: <%v>", conf.logFile, err)
	}
	log.SetOutput(f)

	// customize log format to get date for each line of the log
	log.SetFlags(log.Ldate | log.Lmicroseconds)

	// save pointer to opened log file to close it gracefully when exiting
	conf.logFileHAndle = f

	// build resolver full address
	conf.resolverAddress = fmt.Sprintf("%s:53", conf.resolver)

	// read YAML config
	conf.readBlocklists()

	// var yamlConf YAMLConfig
	// yamlConf.read(conf.yamlConfigFile)
	// fmt.Printf("config=%+v\n", yamlConf)

	// // now read blocklists
	// conf.filters.init()

	// for _, list := range yamlConf.Filters.Blacklist {
	// 	conf.filters.blackList.readFilterFile(list)
	// }
	// for _, list := range yamlConf.Filters.Whitelist {
	// 	conf.filters.whiteList.readFilterFile(list)
	// }

	// nothing entered: show help
	// if len(flag.Args()) == 0 {
	// 	fmt.Print(Usage)
	// 	os.Exit(0)
	// }

	return conf
}

// Read the YAML configuration file
func (yamlConf *YAMLConfig) read(configFile string) *YAMLConfig {
	yamlFile, err := ioutil.ReadFile(configFile)
	if err != nil {
		log.Fatalf("error <%v> opening YAML configuration file: <%s>", err, configFile)
	}
	err = yaml.Unmarshal(yamlFile, yamlConf)
	if err != nil {
		log.Fatalf("error <%v> reading YAML configuration file: <%s>", err, configFile)
	}
	log.Printf("succesfully read YAML file: <%s>, data: <%+v>\n", configFile, yamlConf)

	return yamlConf
}

// Read blocklists and convert them into regexes
func (conf *Config) readBlocklists() {
	//defer conf.mu.Unlock()

	// read YAML config
	var yamlConf YAMLConfig
	yamlConf.read(conf.yamlConfigFile)
	fmt.Printf("config=%+v\n", yamlConf)

	// now read blocklists
	conf.mu.Lock()
	conf.filters.init()

	for _, list := range yamlConf.Filters.Blacklist {
		conf.filters.blackList.readFilterFile(list)
	}
	for _, list := range yamlConf.Filters.Whitelist {
		conf.filters.whiteList.readFilterFile(list)
	}
	conf.mu.Unlock()
}
