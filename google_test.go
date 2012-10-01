package robotstxt

import (
	"testing"
)

const (
	robots_case_order = `user-agent: googlebot-news
Disallow: /
user-agent: *
Disallow: /
user-agent: googlebot
Disallow: /`

	robots_case_grouping = `user-agent: a
disallow: /c

user-agent: b
disallow: /d

user-agent: e
user-agent: f
disallow: /g`

	robots_case_sitemaps = `sitemap: http://test.com/a
user-agent: a
disallow: /c
sitemap: http://test.com/b
user-agent: b
disallow: /d
user-agent: e
sitemap: http://test.com/c
user-agent: f
disallow: /g`

	robots_case_delays = `useragent: a
# some comment : with colon
disallow: /c
user-agent: b
crawldelay: 3.5
disallow: /d
user-agent: e
sitemap: http://test.com/c
user-agent: f
disallow: /g
crawl-delay: 5`

	robots_case_wildcards = `user-agent: *
Disallow: /path*l$`
)

func TestGroupOrder(t *testing.T) {
	agents := []string{"Googlebot-News (Googlebot)", "Googlebot", "Googlebot-Image (Googlebot)", "Otherbot (web)", "Otherbot (News)"}
	groups := []int{1, 3, 3, 2, 2}

	if r, e := FromString(robots_case_order, false); e != nil {
		t.Fatal(e)
	} else {
		for i, a := range agents {
			g := r.findGroup(a)
			gi := getIndexInSlice(r.groups, g) + 1
			if gi != groups[i] {
				t.Fatalf("Expected agent %s to have group number %d, got %d.", a, groups[i], gi)
			}
		}
	}
}

func TestGrouping(t *testing.T) {
	if r, e := FromString(robots_case_grouping, false); e != nil {
		t.Fatal(e)
	} else {
		if len(r.groups) != 3 {
			t.Fatalf("Expected 3 groups, got %d", len(r.groups))
		}
		if len(r.groups[2].agents) != 2 {
			t.Fatalf("Expected 2 agents in group 3, got %d", len(r.groups[2].agents))
		}
		if r.groups[2].agents[0] != "e" {
			t.Fatalf("Expected first agent in group 3 to be e, got %s", len(r.groups[2].agents[0]))
		}
		if r.groups[2].agents[1] != "f" {
			t.Fatalf("Expected second agent in group 3 to be f, got %s", len(r.groups[2].agents[1]))
		}
	}
}

func TestSitemaps(t *testing.T) {
	if r, e := FromString(robots_case_sitemaps, false); e != nil {
		t.Fatal(e)
	} else {
		if len(r.sitemaps) != 3 {
			for i, s := range r.sitemaps {
				t.Logf("Sitemap %d: %s", i, s)
			}
			t.Fatalf("Expected 3 sitemaps, got %d", len(r.sitemaps))
		}
		if len(r.groups) != 3 {
			t.Fatalf("Expected 3 groups, got %d", len(r.groups))
		}
	}
}

func TestCrawlDelays(t *testing.T) {
	if r, e := FromString(robots_case_delays, false); e != nil {
		t.Fatal(e)
	} else {
		if len(r.sitemaps) != 1 {
			t.Fatalf("Expected 1 sitemaps, got %d", len(r.sitemaps))
		}
		if len(r.groups) != 3 {
			t.Fatalf("Expected 3 groups, got %d", len(r.groups))
		}
		if r.groups[1].crawlDelay != 3.5 {
			t.Fatalf("Expected crawl delay of 3.5 for group 2, got %f", r.groups[1].crawlDelay)
		}
		if r.groups[2].crawlDelay != 5 {
			t.Fatalf("Expected crawl delay of 5 for group 3, got %f", r.groups[2].crawlDelay)
		}
	}
}

func TestWildcards(t *testing.T) {
	if r, e := FromString(robots_case_wildcards, false); e != nil {
		t.Fatal(e)
	} else {
		if s := r.groups[0].rules[0].pattern.String(); s != "/path.*l$" {
			t.Fatalf("Expected pattern to be /path.*l$, got %s", s)
		}
	}
}

func getIndexInSlice(ar []*group, g *group) int {
	for i, v := range ar {
		if v == g {
			return i
		}
	}
	return -1
}
