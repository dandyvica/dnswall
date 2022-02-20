package main

import (
	"bufio"
	"fmt"

	//"fmt"
	"log"
	"os"
	"regexp"
	"strings"
)

// List of domains to accept or reject. White list is tested first
type FilteredDomains struct {
	whiteList RegexpFilter
	blackList RegexpFilter
}

// Allocate memory for slice of regexes
func (fd *FilteredDomains) init() {
	fd.whiteList.exprList = make([]*regexp.Regexp, 0)
	fd.blackList.exprList = make([]*regexp.Regexp, 0)
}

// test whether a domain has to be filtered or not
func (domains *FilteredDomains) isFiltered(domain string) bool {
	// try to match a domain in the whitelist first
	for _, re := range domains.whiteList.exprList {
		if re.MatchString(domain) {
			return false
		}
	}

	// try then to match a domain in the blacklist
	for _, re := range domains.blackList.exprList {
		if re.MatchString(domain) {
			fmt.Printf("domain <%s> matched <%s>\n", domain, re.String())
			return true
		}
	}

	return false
}

// When reading a blocklist containing regexes, all data are kept here.
// Each line is converted to a compiled regexp
type RegexpFilter struct {
	exprList []*regexp.Regexp // list of compiled regexes coming from the blocklist
}

// Read a blocklist with one regex per file and create the RegexpFilter struct
// exit process if a regex doesn't compile
func (filter *RegexpFilter) readFilterFile(filterFile string) {
	fileHandle, err := os.Open(filterFile)
	if err != nil {
		log.Fatal(err)
	}
	defer fileHandle.Close()

	scanner := bufio.NewScanner(fileHandle)
	for scanner.Scan() {
		// get rid of trailing spaces
		text := strings.TrimSpace(scanner.Text())

		// skip comments
		if strings.HasPrefix(text, "#") {
			continue
		}

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
