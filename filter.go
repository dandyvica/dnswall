package main

import (
	"regexp"
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
