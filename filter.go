package main

import (
	"bufio"
	//"fmt"
	"log"
	"os"
	"regexp"
	"strings"
)

// List of domains to accept or reject. White list is tested first
type FilteredDomains struct {
	whiteList []regexp.Regexp
	blackList []regexp.Regexp
}

func (domains *FilteredDomains) New() {
	//domains.whiteList
	domains.blackList = append(domains.blackList, *regexp.MustCompile(`\.cn$`))
}

// test whether a domain has to be filtered or not
func (domains *FilteredDomains) IsFiltered(domain string) bool {
	// try to match a domain in the whitelist first
	for _, domainRe := range domains.whiteList {
		if domainRe.MatchString(domain) {
			return false
		}
	}

	// try then to match a domain in the blacklist
	for _, domainRe := range domains.blackList {
		if domainRe.MatchString(domain) {
			return true
		}
	}

	return false
}

// When reading a blocklist containing regexes, all data are kept here.
// Each line is converted to a compiled regexp
type RegexpFilter struct {
	// blocklist file name
	filterFile string

	// list of compiled regexes coming from the blocklist
	exprList []*regexp.Regexp
}

// Read a blocklist with one regex per file and create the RegexpFilter struct
// exit process if a regex doesn't compile
// The initial capacity is made for initializing the slice
func ReadFilterFile(filterFile string, initialCapacity int) *RegexpFilter {
	filter := new(RegexpFilter)
	filter.exprList = make([]*regexp.Regexp, 0, initialCapacity)
	filter.filterFile = filterFile

	fileHandle, err := os.Open(filterFile)
	if err != nil {
		log.Fatal(err)
	}
	defer fileHandle.Close()

	scanner := bufio.NewScanner(fileHandle)
	// optionally, resize scanner's capacity for lines over 64K, see next example
	for scanner.Scan() {
		// get rid of trailing spaces
		text := strings.TrimSpace(scanner.Text())

		// compile the string regexp
		re, err := regexp.Compile(text)
		if err != nil {
			log.Fatalf("regexp <%s> couldn't be compiled, error:<%v>", text, err)
		}

		// add to our list
		filter.exprList = append(filter.exprList, re)
	}

	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}

	return filter
}

// Return true if any of the regexes matches the text
// false otherwise
func (filterList *RegexpFilter) IsMatch(text string) bool {
	for _, expr := range filterList.exprList {
		if expr.MatchString(text) {
			return true
		}
	}
	return false
}
