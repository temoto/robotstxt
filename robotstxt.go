// Package robotstxt implements the robots.txt Exclusion Protocol
// as specified in http://www.robotstxt.org/wc/robots.html
// with various extensions.
package robotstxt

// Comments explaining the logic are taken from either the Google's spec:
// https://developers.google.com/webmasters/control-crawl-index/docs/robots_txt

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"regexp"
	"strconv"
	"strings"
	"time"
)

const (
	AnyGroupId                   = "*"
	regexToRemoveAllPairTagsHTML = `<.*?>.*?</.*?>|<.*?>|<!.*?>`
)

type RobotsData struct {
	// private
	groups      map[string]*Group
	allowAll    bool
	disallowAll bool
	Host        string
	Sitemaps    []string
}

type Group struct {
	rules      []*rule
	Agent      string
	CrawlDelay time.Duration
}

type rule struct {
	path    string
	allow   bool
	pattern *regexp.Regexp
}

type ParseError struct {
	Errs []error
}

func newParseError(errs []error) *ParseError {
	return &ParseError{errs}
}

func (e ParseError) Error() string {
	var b bytes.Buffer

	b.WriteString("Parse error(s): " + "\n")
	for _, er := range e.Errs {
		b.WriteString(er.Error() + "\n")
	}
	return b.String()
}

var allowAll = &RobotsData{allowAll: true}
var disallowAll = &RobotsData{disallowAll: true}
var emptyGroup = &Group{}

func FromStatusAndBytes(statusCode int, body []byte) (*RobotsData, error) {
	switch {
	case statusCode >= 200 && statusCode < 300:
		return FromBytes(body)

	// From https://developers.google.com/webmasters/control-crawl-index/docs/robots_txt
	//
	// Google treats all 4xx errors in the same way and assumes that no valid
	// robots.txt file exists. It is assumed that there are no restrictions.
	// This is a "full allow" for crawling. Note: this includes 401
	// "Unauthorized" and 403 "Forbidden" HTTP result codes.
	case statusCode >= 400 && statusCode < 500:
		return allowAll, nil

	// From Google's spec:
	// Server errors (5xx) are seen as temporary errors that result in a "full
	// disallow" of crawling.
	case statusCode >= 500 && statusCode < 600:
		return disallowAll, nil
	}

	return nil, errors.New("Unexpected status: " + strconv.Itoa(statusCode))
}

func FromStatusAndString(statusCode int, body string) (*RobotsData, error) {
	return FromStatusAndBytes(statusCode, []byte(body))
}

func FromResponse(res *http.Response) (*RobotsData, error) {
	if res == nil {
		// Edge case, if res is nil, return nil data
		return nil, nil
	}
	buf, e := ioutil.ReadAll(res.Body)
	if e != nil {
		return nil, e
	}
	return FromStatusAndBytes(res.StatusCode, buf)
}

// This method uses a regular expresion to remove HTML.
// must be of the form <tag> txt <\tag> ie well formed.
func stripHtmlRegex(s string) string {
	r1 := regexp.MustCompile(regexToRemoveAllPairTagsHTML)
	return r1.ReplaceAllString(s, "")
}

func FromBytes(body []byte) (r *RobotsData, err error) {
	var errs []error

	// special case (probably not worth optimization?)
	trimmed := bytes.TrimSpace(body)
	if len(trimmed) == 0 {
		return allowAll, nil
	}

	// toss any html
	trimedFromHtml := stripHtmlRegex(string(body))
	if len(trimedFromHtml) == 0 {
		return allowAll, nil
	}
	body = []byte(trimedFromHtml)

	sc := newByteScanner("bytes", true)
	// sc.Quiet = !print_errors
	sc.feed(body, true)
	tokens := sc.scanAll()

	// special case worth optimization
	if len(tokens) == 0 {
		return allowAll, nil
	}

	r = &RobotsData{}
	parser := newParser(tokens)
	r.groups, r.Host, r.Sitemaps, errs = parser.parseAll()
	if len(errs) > 0 {
		return nil, newParseError(errs)
	}

	return r, nil
}

func FromString(body string) (r *RobotsData, err error) {
	return FromBytes([]byte(body))
}

func (r *RobotsData) TestAgent(path, agent string) bool {
	if r.allowAll {
		return true
	}
	if r.disallowAll {
		return false
	}

	// Find a group of rules that applies to this agent
	// From Google's spec:
	// The user-agent is non-case-sensitive.
	g := r.FindGroup(agent)
	return g.Test(path)
}

func (r *RobotsData) TestGroup(path string, group *Group) bool {
	if r.allowAll {
		return true
	}
	if r.disallowAll {
		return false
	}

	return group.Test(path)
}

// Returns true if all urls disallowed
func (r *RobotsData) TestDisallowAll() bool {
	return !r.disallowAll
}

// FindGroup searches block of declarations for specified user-agent.
// From Google's spec:
// Only one group of group-member records is valid for a particular crawler.
// The crawler must determine the correct group of records by finding the group
// with the most specific user-agent that still matches. All other groups of
// records are ignored by the crawler. The user-agent is non-case-sensitive.
// The order of the groups within the robots.txt file is irrelevant.
func (r *RobotsData) FindGroup(agent string) (ret *Group) {
	_, g := r.FindGroupWithGroupId(agent)
	return g
}

func (r *RobotsData) FindGroupWithGroupId(agent string) (groupId string, ret *Group) {
	var prefixLen int

	agent = strings.ToLower(agent)
	if ret = r.groups[AnyGroupId]; ret != nil {
		// Weakest match possible
		prefixLen = 1
		groupId = AnyGroupId
	}
	for a, g := range r.groups {
		if a != AnyGroupId && strings.HasPrefix(agent, a) {
			if l := len(a); l > prefixLen {
				prefixLen = l
				ret = g
				groupId = a
			}
		}
	}

	if ret == nil {
		return AnyGroupId, emptyGroup
	}
	return
}

func (r *RobotsData) SetGroups(groups map[string]*Group) {
	r.groups = groups
}

func(r *RobotsData) MarshalJSON() ([]byte, error) {
	return json.Marshal(map[string]interface{}{
		"allow_all":r.allowAll,
		"disallow_all":r.disallowAll,
		"groups":r.groups,
		"host":r.Host,
		"sitemaps":r.Sitemaps,
	})
}

func(r *RobotsData) UnmarshalJSON(bytes []byte) error {
	var robotsDataInterface map[string]interface{}
	err := json.Unmarshal(bytes, &robotsDataInterface)

	if err != nil {
		return err
	}

	if allowAll, ok := robotsDataInterface["allow_all"].(bool); ok {
		r.allowAll = allowAll
	}

	if disallowAll, ok := robotsDataInterface["disallow_all"].(bool); ok {
		r.disallowAll = disallowAll
	}

	if groupInterfaces, ok := robotsDataInterface["groups"].(map[string]interface{}); ok {

		r.groups = make(map[string]*Group, len(groupInterfaces))

		for key, groupInterface := range groupInterfaces {
			g := Group{}
			err = groupInterfaceToGroup(groupInterface, &g)

			if err != nil {
				return err
			}

			r.groups[key] = &g
		}
	}

	if host, ok := robotsDataInterface["host"].(string); ok {
		r.Host = host
	}

	if sitemaps, ok := robotsDataInterface["sitemaps"].([]string); ok {
		r.Sitemaps = sitemaps
	}

	return nil
}

func groupInterfaceToGroup(groupInterface interface{}, group *Group) error {
	groupMapInterface, ok := groupInterface.(map[string]interface{})

	if !ok {
		return fmt.Errorf("Could not parse Group interface")
	}

	if agent, ok := groupMapInterface["agent"].(string); ok {
		group.Agent = agent
	}

	if crawlDelay, ok := groupMapInterface["crawl_delay"].(float64); ok {
		group.CrawlDelay = time.Duration(int64(crawlDelay))
	}

	if rulesAr, ok := groupMapInterface["rules"].([]interface{}); ok && len(rulesAr) > 0 {

		group.rules = make([]*rule, 0, len(rulesAr))

		for _, ruleInterface := range rulesAr {
			restoredRule := rule{}
			err := ruleInterfaceToRule(ruleInterface, &restoredRule)

			if err != nil {
				return err
			}

			group.rules = append(group.rules, &restoredRule)
		}
	}

	return nil
}

func ruleInterfaceToRule(ruleInterface interface{}, restoredRule *rule) error {
	r, ok := ruleInterface.(map[string]interface{})

	if !ok {
		return fmt.Errorf("Could not parse Rule Interface")
	}

	if allow, ok := r["allow"].(bool); ok {
		restoredRule.allow = allow
	}

	if path, ok := r["path"].(string); ok {
		restoredRule.path = path
	}

	if pattern, ok := r["pattern"].(string); ok && len(pattern) > 0 {
		var err error
		restoredRule.pattern, err = regexp.Compile(pattern)

		if err != nil {
			return err
		}
	}

	return nil
}

func (g *Group) Test(path string) bool {
	if r := g.findRule(path); r != nil {
		return r.allow
	}

	// From Google's spec:
	// By default, there are no restrictions for crawling for the designated crawlers.
	return true
}

// From Google's spec:
// The path value is used as a basis to determine whether or not a rule applies
// to a specific URL on a site. With the exception of wildcards, the path is
// used to match the beginning of a URL (and any valid URLs that start with the
// same path).
//
// At a group-member level, in particular for allow and disallow directives,
// the most specific rule based on the length of the [path] entry will trump
// the less specific (shorter) rule. The order of precedence for rules with
// wildcards is undefined.
func (g *Group) findRule(path string) (ret *rule) {
	var prefixLen int

	for _, r := range g.rules {
		if r.pattern != nil {
			if r.pattern.MatchString(path) {
				// Consider this a match equal to the length of the pattern.
				// From Google's spec:
				// The order of precedence for rules with wildcards is undefined.
				if l := len(r.pattern.String()); l > prefixLen {
					prefixLen = l
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

func(g *Group) MarshalJSON() ([]byte, error) {
	return json.Marshal(map[string]interface{}{
		"agent":g.Agent,
		"crawl_delay":g.CrawlDelay.Nanoseconds(),
		"rules":g.rules,
	})
}

func(r *rule) MarshalJSON() ([]byte, error) {
	var pattern string

	if r.pattern != nil {
		pattern = r.pattern.String()
	}

	return json.Marshal(map[string]interface{}{
		"allow":r.allow,
		"path":r.path,
		"pattern":pattern,
	})
}

