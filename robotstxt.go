// The robots.txt Exclusion Protocol is implemented as specified in
// http://www.robotstxt.org/wc/robots.html
// with various extensions.
package robotstxt

import (
	"bytes"
	"errors"
	"strings"
)

type RobotsData struct {
	DefaultAgent string
	// private
	groups      []group
	allowAll    bool
	disallowAll bool
	sitemaps    []string
}

type group struct {
	agent      string
	rules      []rule
	crawlDelay uint
}

type rule struct {
	path  string
	allow bool
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
		return AllowAll, nil
	case statusCode >= 200 && statusCode < 300:
		return FromBytes(body, print_errors)
	}
	// Conservative disallow all default
	//
	// From https://developers.google.com/webmasters/control-crawl-index/docs/robots_txt
	//
	// Server errors (5xx) are seen as temporary errors that result in a "full
	// disallow" of crawling.
	return DisallowAll, nil
}

func FromResponse(statusCode int, body string, print_errors bool) (*RobotsData, error) {
	return FromResponseBytes(statusCode, []byte(body), print_errors)
}

func FromBytes(body []byte, print_errors bool) (r *RobotsData, err error) {
	// special case (probably not worth optimization?)
	trimmed := bytes.TrimSpace(body)
	if len(trimmed) == 0 {
		return AllowAll, nil
	}

	sc := NewByteScanner("bytes", false)
	sc.Quiet = !print_errors
	sc.Feed(body, true)
	var tokens []string
	tokens, err = sc.ScanAll()
	if err != nil {
		return nil, err
	}

	// special case worth optimization
	if len(tokens) == 0 {
		return AllowAll, nil
	}

	r = &RobotsData{}
	parser := NewParser(tokens)
	r.rules, err = parser.ParseAll()

	return r, err
}

func FromString(body string, print_errors bool) (r *RobotsData, err error) {
	return FromBytes([]byte(body), print_errors)
}

func (r *RobotsData) Test(url string) bool {
	return r.TestAgent(url, r.DefaultAgent)
}

func (r *RobotsData) TestAgent(url, agent string) (allow bool) {
	if r.allowAll {
		return true
	}
	if r.disallowAll {
		return false
	}

	// optimistic
	allow = true
	for _, rule := range r.rules {
		if rule.MatchAgent(agent) && rule.MatchUrl(url) {
			allow = rule.Allow
			// stop on first disallow as safety default
			// in absense of better algorithm
			if !rule.Allow {
				break
			}
		}
	}

	return allow
}

func (rule *Rule) MatchAgent(agent string) bool {
	l_agent := strings.ToLower(agent)
	l_rule_agent := strings.ToLower(rule.Agent)
	return rule.Agent == "*" || strings.HasPrefix(l_agent, l_rule_agent)
}

func (rule *Rule) MatchUrl(url string) bool {
	return strings.HasPrefix(url, rule.Uri)
}

func (rule *Rule) String() string {
	allow_str := "Disallow"
	if rule.Allow {
		allow_str = "Allow"
	}
	return "<" + allow_str + " " + rule.Agent + " " + rule.Uri + ">"
}
