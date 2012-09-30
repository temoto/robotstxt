package robotstxt

// Comments explaining the logic are taken from either the google's spec:
// https://developers.google.com/webmasters/control-crawl-index/docs/robots_txt
//
// or the Wikipedia's entry on robots.txt:
// http://en.wikipedia.org/wiki/Robots.txt

import (
	"errors"
	"fmt"
	"io"
	"regexp"
	"strconv"
	"strings"
)

type lineType uint

const (
	lIgnore lineType = iota
	lUnknown
	lUserAgent
	lAllow
	lDisallow
	lCrawlDelay
	lSitemap
)

type parser struct {
	tokens []string
	pos    int
}

type lineInfo struct {
	t  lineType
	k  string
	vs string
	vf float64
	vr *regexp.Regexp
}

func newParser(tokens []string) *parser {
	return &parser{tokens: tokens}
}

func (p *parser) parseAll() (groups []*group, sitemaps []string, errs []error) {
	var curGroup *group

	// Reset internal fields, tokens are assigned at creation time, never change
	p.pos = 0

	// TODO : Two successive user-agent lines are part of the same group, so a group
	// may apply to more than one user-agent!
	// Re: Google's spec:
	// There are three distinct groups specified, one for "a" and one for "b"
	// as well as one for both "e" and "f".

	for {
		if li, err := p.parseLine(); err != nil {
			if err == io.EOF {
				// Append the current group if any
				if curGroup != nil {
					groups = append(groups, curGroup)
				}
				break
			}
			errs = append(errs, err)
		} else {
			switch li.t {
			case lUserAgent:
				// End previous group
				if curGroup != nil {
					groups = append(groups, curGroup)
				}
				// Start new group
				curGroup = &group{agent: li.vs}
			case lDisallow:
				// Error if no current group
				if curGroup == nil {
					errs = append(errs, errors.New(fmt.Sprintf("Disallow before User-agent at token #%d.", p.pos)))
				} else {
					curGroup.rules = append(curGroup.rules, &rule{li.vs, false, nil})
				}
			case lAllow:
				// Error if no current group
				if curGroup == nil {
					errs = append(errs, errors.New(fmt.Sprintf("Allow before User-agent at token #%d.", p.pos)))
				} else {
					curGroup.rules = append(curGroup.rules, &rule{li.vs, true, nil})
				}
			case lSitemap:
				sitemaps = append(sitemaps, li.vs)
			case lCrawlDelay:
				if curGroup == nil {
					errs = append(errs, errors.New(fmt.Sprintf("Crawl-delay before User-agent at token #%d.", p.pos)))
				} else {
					curGroup.crawlDelay = li.vf
				}
			}
		}
	}
	if len(errs) > 0 {
		for _, e := range errs {
			fmt.Printf("Error: %s\n", e.Error())
		}
	}
	return
}

func (p *parser) parseLine() (li *lineInfo, err error) {
	t1, ok1 := p.popToken()
	if !ok1 {
		// proper EOF
		return nil, io.EOF
	}

	t2, ok2 := p.peekToken()
	if !ok2 {
		// EOF, no value associated with the token, so ignore token and return
		if strings.Trim(t1, " \t\v\n\r") != "" {
			return nil, errors.New(fmt.Sprintf(`Unexpected EOF at token #%d namely: "%s".`, p.pos, t1))
		} else {
			return nil, io.EOF
		}
	}

	// Helper closure for all string-based tokens, common behaviour:
	// - Consume t2 token
	// - If empty, return unkown line info
	// - Otherwise return the specified line info
	returnStringVal := func(t lineType) (*lineInfo, error) {
		p.popToken()
		if t2 != "" {
			return &lineInfo{t: t, k: t1, vs: t2}, nil
		}
		return &lineInfo{t: lIgnore}, nil
	}

	// TODO : For paths, automatically add the starting "/", ignore a trailing "*",
	// and manage wildcards within a path (turn into a pattern)

	switch strings.ToLower(t1) {
	case "\n":
		// Don't consume t2 and continue parsing
		return &lineInfo{t: lIgnore}, nil

	case "user-agent", "useragent":
		// From google's spec:
		// Handling of <field> elements with simple errors / typos (eg "useragent"
		// instead of "user-agent") is undefined and may be interpreted as correct
		// directives by some user-agents.
		return returnStringVal(lUserAgent)

	case "disallow":
		// From google's spec:
		// When no path is specified, the directive is ignored (so an empty Disallow
		// CAN be an allow, since allow is the default. The actual result depends
		// on the other rules in the group).
		return returnStringVal(lDisallow)

	case "allow":
		// From google's spec:
		// When no path is specified, the directive is ignored.
		return returnStringVal(lAllow)

	case "sitemap":
		// Non-group field, applies to the host as a whole, not to a specific user-agent
		return returnStringVal(lSitemap)

	case "crawl-delay", "crawldelay":
		// From http://en.wikipedia.org/wiki/Robots_exclusion_standard#Nonstandard_extensions
		// Several major crawlers support a Crawl-delay parameter, set to the
		// number of seconds to wait between successive requests to the same server.
		p.popToken()
		if cd, e := strconv.ParseFloat(t2, 64); e != nil {
			return nil, e
		} else {
			return &lineInfo{t: lCrawlDelay, k: t1, vf: cd}, nil
		}
	}

	// Consume t2 token
	//p.popToken()
	return &lineInfo{t: lUnknown, k: t1}, nil
}

func (p *parser) popToken() (tok string, ok bool) {
	tok, ok = p.peekToken()
	if !ok {
		return
	}
	p.pos++
	return tok, true
}

func (p *parser) peekToken() (tok string, ok bool) {
	if p.pos >= len(p.tokens) {
		return "", false
	}
	return p.tokens[p.pos], true
}
