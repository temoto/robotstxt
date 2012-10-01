// The robots.txt Exclusion Protocol is implemented as specified in
// http://www.robotstxt.org/wc/robots.html
// with various extensions.
package robotstxt

// Comments explaining the logic are taken from either the google's spec:
// https://developers.google.com/webmasters/control-crawl-index/docs/robots_txt

import (
	"bytes"
	"errors"
	"fmt"
	"os"
	"regexp"
	"strings"
)

type RobotsData struct {
	DefaultAgent string
	// private
	groups      []*group
	allowAll    bool
	disallowAll bool
	sitemaps    []string
}

type group struct {
	agents     []string
	rules      []*rule
	crawlDelay float64
}

type rule struct {
	path    string
	allow   bool
	pattern *regexp.Regexp
}

var allowAll = &RobotsData{allowAll: true}
var disallowAll = &RobotsData{disallowAll: true}

func FromResponseBytes(statusCode int, body []byte, print_errors bool) (*RobotsData, error) {
	switch {
	//
	// From https://developers.google.com/webmasters/control-crawl-index/docs/robots_txt
	//
	// Google treats all 4xx errors in the same way and assumes that no valid
	// robots.txt file exists. It is assumed that there are no restrictions.
	// This is a "full allow" for crawling. Note: this includes 401
	// "Unauthorized" and 403 "Forbidden" HTTP result codes.
	case statusCode >= 400 && statusCode < 500:
		return allowAll, nil
	case statusCode >= 200 && statusCode < 300:
		return FromBytes(body, print_errors)
	}
	// Conservative disallow all default
	//
	// From https://developers.google.com/webmasters/control-crawl-index/docs/robots_txt
	//
	// Server errors (5xx) are seen as temporary errors that result in a "full
	// disallow" of crawling.
	return disallowAll, nil
}

func FromResponse(statusCode int, body string, print_errors bool) (*RobotsData, error) {
	return FromResponseBytes(statusCode, []byte(body), print_errors)
}

func FromBytes(body []byte, print_errors bool) (r *RobotsData, err error) {
	var errs []error

	// special case (probably not worth optimization?)
	trimmed := bytes.TrimSpace(body)
	if len(trimmed) == 0 {
		return allowAll, nil
	}

	sc := newByteScanner("bytes", false)
	sc.Quiet = !print_errors
	sc.Feed(body, true)
	var tokens []string
	tokens, err = sc.ScanAll()
	if err != nil {
		return nil, err
	}

	// special case worth optimization
	if len(tokens) == 0 {
		return allowAll, nil
	}

	r = &RobotsData{}
	parser := newParser(tokens)
	r.groups, r.sitemaps, errs = parser.parseAll()
	if len(errs) > 0 {
		if print_errors {
			for _, e := range errs {
				fmt.Fprintln(os.Stderr, e)
			}
		}
		return nil, errors.New("Parse error. Use print_errors = true to print on stderr.")
	}

	return r, nil
}

func FromString(body string, print_errors bool) (r *RobotsData, err error) {
	return FromBytes([]byte(body), print_errors)
}

func (r *RobotsData) Test(path string) bool {
	return r.TestAgent(path, r.DefaultAgent)
}

func (r *RobotsData) TestAgent(path, agent string) bool {
	if r.allowAll {
		return true
	}
	if r.disallowAll {
		return false
	}

	// Find a group of rules that applies to this agent
	// From google's spec:
	// The user-agent is non-case-sensitive.
	if g := r.findGroup(agent); g != nil {
		// Find a rule that applies to this url
		if r := g.findRule(path); r != nil {
			return r.allow
		}
	}

	// From google's spec:
	// By default, there are no restrictions for crawling for the designated crawlers. 
	return true
}

// From google's spec:
// Only one group of group-member records is valid for a particular crawler.
// The crawler must determine the correct group of records by finding the group
// with the most specific user-agent that still matches. All other groups of
// records are ignored by the crawler. The user-agent is non-case-sensitive.
// The order of the groups within the robots.txt file is irrelevant.
func (r *RobotsData) findGroup(agent string) (ret *group) {
	var prefixLen int

	agent = strings.ToLower(agent)
	for _, g := range r.groups {
		for _, a := range g.agents {
			if a == "*" && prefixLen == 0 {
				// Weakest match possible
				prefixLen = 1
				ret = g
			} else if strings.HasPrefix(agent, a) {
				if l := len(a); l > prefixLen {
					prefixLen = l
					ret = g
				}
			}
		}
	}
	return
}

// From google's spec:
// The path value is used as a basis to determine whether or not a rule applies
// to a specific URL on a site. With the exception of wildcards, the path is
// used to match the beginning of a URL (and any valid URLs that start with the
// same path).
//
// At a group-member level, in particular for allow and disallow directives,
// the most specific rule based on the length of the [path] entry will trump
// the less specific (shorter) rule. The order of precedence for rules with
// wildcards is undefined.
func (g *group) findRule(path string) (ret *rule) {
	var prefixLen int

	for _, r := range g.rules {
		if r.pattern != nil {
			if r.pattern.MatchString(path) {
				// Consider this a match equal to the length of the pattern.
				// From google's spec:
				// The order of precedence for rules with wildcards is undefined.
				if l := len(r.pattern.String()); l > prefixLen {
					prefixLen = len(r.pattern.String())
					ret = r
				}
			}
		} else if r.path == "/" && prefixLen == 0 {
			// Weakest match possible
			prefixLen = 1
			ret = r
		} else if strings.HasPrefix(path, r.path) {
			if l := len(r.path); l > prefixLen {
				prefixLen = l
				ret = r
			}
		}
	}
	return
}
