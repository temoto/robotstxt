package robotstxt

import (
	"testing"
)

func TestFromResponseBasic(t *testing.T) {
	if _, err := FromResponse(200, "", true); err != nil {
		t.Fatal("FromResponse MUST accept 200/\"\"")
	}
	if _, err := FromResponse(401, "", true); err != nil {
		t.Fatal("FromResponse MUST accept 401/\"\"")
	}
	if _, err := FromResponse(403, "", true); err != nil {
		t.Fatal("FromResponse MUST accept 403/\"\"")
	}
	if _, err := FromResponse(404, "", true); err != nil {
		t.Fatal("FromResponse MUST accept 404/\"\"")
	}
}

func _expectAllow(r *RobotsData, t *testing.T) bool {
	allow := r.TestAgent("/", "Somebot")
	return allow
}

func ExpectAllow(r *RobotsData, t *testing.T, msg string) {
	if !_expectAllow(r, t) {
		t.Fatal(msg)
	}
}

func ExpectDisallow(r *RobotsData, t *testing.T, msg string) {
	if _expectAllow(r, t) {
		t.Fatal(msg)
	}
}

func TestResponse401(t *testing.T) {
	r, _ := FromResponse(401, "", true)
	ExpectAllow(r, t, "FromResponse(401, \"\") MUST allow everything.")
}

func TestResponse403(t *testing.T) {
	r, _ := FromResponse(403, "", true)
	ExpectAllow(r, t, "FromResponse(403, \"\") MUST allow everything.")
}

func TestResponse404(t *testing.T) {
	r, _ := FromResponse(404, "", true)
	ExpectAllow(r, t, "FromResponse(404, \"\") MUST allow everything.")
}

func TestFromStringBasic(t *testing.T) {
	if _, err := FromString("", true); err != nil {
		t.Fatal("FromString MUST accept \"\"")
	}
}

func TestFromStringEmpty(t *testing.T) {
	r, _ := FromString("", true)
	if allow := r.TestAgent("/", "Somebot"); !allow {
		t.Fatal("FromString(\"\") MUST allow everything.")
	}
}

func TestFromStringComment(t *testing.T) {
	if _, err := FromString("# comment", true); err != nil {
		t.Fatal("FromString MUST accept \"# comment\"")
	}
}

func TestFromString001(t *testing.T) {
	r, err := FromString("User-Agent: *\r\nDisallow: /\r\n", true)
	if err != nil {
		t.Fatal(err.Error())
	}
	allow := r.TestAgent("/foobar", "SomeAgent")
	if allow {
		t.Fatal("Must deny.")
	}
}

func TestFromString002(t *testing.T) {
	r, err := FromString("User-Agent: *\r\nDisallow: /account\r\n", true)
	if err != nil {
		t.Fatal(err.Error())
	}
	allow := r.TestAgent("/foobar", "SomeAgent")
	if !allow {
		t.Fatal("Must allow.")
	}
}

const robots_text_001 = "User-agent: * \nDisallow: /administrator/\nDisallow: /cache/\nDisallow: /components/\nDisallow: /editor/\nDisallow: /forum/\nDisallow: /help/\nDisallow: /images/\nDisallow: /includes/\nDisallow: /language/\nDisallow: /mambots/\nDisallow: /media/\nDisallow: /modules/\nDisallow: /templates/\nDisallow: /installation/\nDisallow: /getcid/\nDisallow: /tooltip/\nDisallow: /getuser/\nDisallow: /download/\nDisallow: /index.php?option=com_phorum*,quote=1\nDisallow: /index.php?option=com_phorum*phorum_query=search\nDisallow: /index.php?option=com_phorum*,newer\nDisallow: /index.php?option=com_phorum*,older\n\nUser-agent: Yandex\nAllow: /\nSitemap: http://www.pravorulya.com/sitemap.xml\nSitemap: http://www.pravorulya.com/sitemap1.xml"

func TestFromString003(t *testing.T) {
	r, err := FromString(robots_text_001, true)
	if err != nil {
		t.Fatal(err.Error())
	}
	allow := r.TestAgent("/administrator/", "SomeBot")
	if allow {
		t.Fatal("Must deny.")
	}
}

func TestFromString004(t *testing.T) {
	r, err := FromString(robots_text_001, true)
	if err != nil {
		t.Fatal(err.Error())
	}
	allow := r.TestAgent("/paruram", "SomeBot")
	if !allow {
		t.Fatal("Must allow.")
	}
}

func TestInvalidEncoding(t *testing.T) {
	// Invalid UTF-8 encoding should not break parser.
	_, err := FromString("User-agent: H\xef\xbf\xbdm�h�kki\nDisallow: *", true)
	if err != nil {
		t.Fatal(err.Error())
	}
}

// http://www.google.com/robots.txt on Wed, 12 Jan 2011 12:22:20 GMT
const robots_text_002 = ("User-agent: *\nDisallow: /search\nDisallow: /groups\nDisallow: /images\nDisallow: /catalogs\nDisallow: /catalogues\nDisallow: /news\nAllow: /news/directory\nDisallow: /nwshp\nDisallow: /setnewsprefs?\nDisallow: /index.html?\nDisallow: /?\nDisallow: /addurl/image?\nDisallow: /pagead/\nDisallow: /relpage/\nDisallow: /relcontent\nDisallow: /imgres\nDisallow: /imglanding\nDisallow: /keyword/\nDisallow: /u/\nDisallow: /univ/\nDisallow: /cobrand\nDisallow: /custom\nDisallow: /advanced_group_search\nDisallow: /googlesite\nDisallow: /preferences\nDisallow: /setprefs\nDisallow: /swr\nDisallow: /url\nDisallow: /default\nDisallow: /m?\nDisallow: /m/?\nDisallow: /m/blogs?\nDisallow: /m/directions?\nDisallow: /m/ig\nDisallow: /m/images?\nDisallow: /m/local?\nDisallow: /m/movies?\nDisallow: /m/news?\nDisallow: /m/news/i?\nDisallow: /m/place?\nDisallow: /m/products?\nDisallow: /m/products/\nDisallow: /m/setnewsprefs?\nDisallow: /m/search?\nDisallow: /m/swmloptin?\nDisallow: /m/trends\nDisallow: /m/video?\nDisallow: /wml?\nDisallow: /wml/?\nDisallow: /wml/search?\nDisallow: /xhtml?\nDisallow: /xhtml/?\nDisallow: /xhtml/search?\nDisallow: /xml?\nDisallow: /imode?\nDisallow: /imode/?\nDisallow: /imode/search?\nDisallow: /jsky?\nDisallow: /jsky/?\nDisallow: /jsky/search?\nDisallow: /pda?\nDisallow: /pda/?\nDisallow: /pda/search?\nDisallow: /sprint_xhtml\nDisallow: /sprint_wml\nDisallow: /pqa\nDisallow: /palm\nDisallow: /gwt/\nDisallow: /purchases\nDisallow: /hws\nDisallow: /bsd?\nDisallow: /linux?\nDisallow: /mac?\nDisallow: /microsoft?\nDisallow: /unclesam?\nDisallow: /answers/search?q=\nDisallow: /local?\nDisallow: /local_url\nDisallow: /froogle?\nDisallow: /products?\nDisallow: /products/\nDisallow: /froogle_\nDisallow: /product_\nDisallow: /products_\nDisallow: /products;\nDisallow: /print\nDisallow: /books\nDisallow: /bkshp?q=\nAllow: /booksrightsholders\nDisallow: /patents?\nDisallow: /patents/\nAllow: /patents/about\nDisallow: /scholar\nDisallow: /complete\nDisallow: /sponsoredlinks\nDisallow: /videosearch?\nDisallow: /videopreview?\nDisallow: /videoprograminfo?\nDisallow: /maps?\nDisallow: /mapstt?\nDisallow: /mapslt?\nDisallow: /maps/stk/\nDisallow: /maps/br?\nDisallow: /mapabcpoi?\nDisallow: /maphp?\nDisallow: /places/\nAllow: /places/$\nDisallow: /maps/place\nDisallow: /help/maps/streetview/partners/welcome/\nDisallow: /lochp?\nDisallow: /center\nDisallow: /ie?\nDisallow: /sms/demo?\nDisallow: /katrina?\nDisallow: /blogsearch?\nDisallow: /blogsearch/\nDisallow: /blogsearch_feeds\nDisallow: /advanced_blog_search\nDisallow: /reader/\nAllow: /reader/play\nDisallow: /uds/\nDisallow: /chart?\nDisallow: /transit?\nDisallow: /mbd?\nDisallow: /extern_js/\nDisallow: /calendar/feeds/\nDisallow: /calendar/ical/\nDisallow: /cl2/feeds/\n" +
	"Disallow: /cl2/ical/\nDisallow: /coop/directory\nDisallow: /coop/manage\nDisallow: /trends?\nDisallow: /trends/music?\nDisallow: /trends/hottrends?\nDisallow: /trends/viz?\nDisallow: /notebook/search?\nDisallow: /musica\nDisallow: /musicad\nDisallow: /musicas\nDisallow: /musicl\nDisallow: /musics\nDisallow: /musicsearch\nDisallow: /musicsp\nDisallow: /musiclp\nDisallow: /browsersync\nDisallow: /call\nDisallow: /archivesearch?\nDisallow: /archivesearch/url\nDisallow: /archivesearch/advanced_search\nDisallow: /base/reportbadoffer\nDisallow: /urchin_test/\nDisallow: /movies?\nDisallow: /codesearch?\nDisallow: /codesearch/feeds/search?\nDisallow: /wapsearch?\nDisallow: /safebrowsing\nAllow: /safebrowsing/diagnostic\nAllow: /safebrowsing/report_error/\nAllow: /safebrowsing/report_phish/\nDisallow: /reviews/search?\nDisallow: /orkut/albums\nAllow: /jsapi\nDisallow: /views?\nDisallow: /c/\nDisallow: /cbk\nDisallow: /recharge/dashboard/car\nDisallow: /recharge/dashboard/static/\nDisallow: /translate_a/\nDisallow: /translate_c\nDisallow: /translate_f\nDisallow: /translate_static/\nDisallow: /translate_suggestion\nDisallow: /profiles/me\nAllow: /profiles\nDisallow: /s2/profiles/me\nAllow: /s2/profiles\nAllow: /s2/photos\nAllow: /s2/static\nDisallow: /s2\nDisallow: /transconsole/portal/\nDisallow: /gcc/\nDisallow: /aclk\nDisallow: /cse?\nDisallow: /cse/home\nDisallow: /cse/panel\nDisallow: /cse/manage\nDisallow: /tbproxy/\nDisallow: /imesync/\nDisallow: /shenghuo/search?\nDisallow: /support/forum/search?\nDisallow: /reviews/polls/\nDisallow: /hosted/images/\nDisallow: /ppob/?\nDisallow: /ppob?\nDisallow: /ig/add?\nDisallow: /adwordsresellers\nDisallow: /accounts/o8\nAllow: /accounts/o8/id\nDisallow: /topicsearch?q=\nDisallow: /xfx7/\nDisallow: /squared/api\nDisallow: /squared/search\nDisallow: /squared/table\nDisallow: /toolkit/\nAllow: /toolkit/*.html\nDisallow: /globalmarketfinder/\nAllow: /globalmarketfinder/*.html\nDisallow: /qnasearch?\nDisallow: /errors/\nDisallow: /app/updates\nDisallow: /sidewiki/entry/\nDisallow: /quality_form?\nDisallow: /labs/popgadget/search\nDisallow: /buzz/post\nDisallow: /compressiontest/\nDisallow: /analytics/reporting/\nDisallow: /analytics/admin/\nDisallow: /analytics/web/\nDisallow: /analytics/feeds/\nDisallow: /analytics/settings/\nDisallow: /alerts/\nDisallow: /phone/compare/?\nAllow: /alerts/manage\nSitemap: http://www.gstatic.com/s2/sitemaps/profiles-sitemap.xml\nSitemap: http://www.google.com/hostednews/sitemap_index.xml\nSitemap: http://www.google.com/ventures/sitemap_ventures.xml\nSitemap: http://www.google.com/sitemaps_webmasters.xml\nSitemap: http://www.gstatic.com/trends/websites/sitemaps/sitemapindex.xml\nSitemap: http://www.gstatic.com/dictionary/static/sitemaps/sitemap_index.xml")

func TestFromString005(t *testing.T) {
	r, err := FromString(robots_text_002, true)
	if err != nil {
		t.Fatal(err.Error())
	}
	ExpectAllow(r, t, "Must allow.")
}

func TestFromString006(t *testing.T) {
	r, err := FromString(robots_text_002, true)
	if err != nil {
		t.Fatal(err.Error())
	}
	allow := r.TestAgent("/search", "SomeBot")
	if allow {
		t.Fatal("Must deny.")
	}
}

const robots_text_003 = "User-Agent: * \nAllow: /"

func TestFromString007(t *testing.T) {
	r, err := FromString(robots_text_003, true)
	if err != nil {
		t.Fatal(err.Error())
	}
	allow := r.TestAgent("/random", "SomeBot")
	if !allow {
		t.Fatal("Must allow.")
	}
}

const robots_text_004 = "User-Agent: * \nDisallow: "

func TestFromString008(t *testing.T) {
	r, err := FromString(robots_text_004, true)
	if err != nil {
		t.Log(robots_text_004)
		t.Fatal(err.Error())
	}
	allow := r.TestAgent("/random", "SomeBot")
	if !allow {
		t.Fatal("Must allow.")
	}
}

const robots_text_005 = `User-agent: Google
Disallow:
User-agent: *
Disallow: /`

func TestRobotstxtOrgCase1(t *testing.T) {
	if r, err := FromString(robots_text_005, false); err != nil {
		t.Fatal(err.Error())
	} else if allow := r.TestAgent("/path/page1.html", "SomeBot"); allow {
		t.Fatal("Must disallow.")
	}
}

func TestRobotstxtOrgCase2(t *testing.T) {
	if r, err := FromString(robots_text_005, false); err != nil {
		t.Fatal(err.Error())
	} else if allow := r.TestAgent("/path/page1.html", "Googlebot"); !allow {
		t.Fatal("Must allow.")
	}
}

func BenchmarkParseFromString001(b *testing.B) {
	for i := 0; i < b.N; i++ {
		FromString(robots_text_001, false)
		b.SetBytes(int64(len(robots_text_001)))
	}
}

func BenchmarkParseFromString002(b *testing.B) {
	for i := 0; i < b.N; i++ {
		FromString(robots_text_002, false)
		b.SetBytes(int64(len(robots_text_002)))
	}
}

func BenchmarkParseFromResponse401(b *testing.B) {
	for i := 0; i < b.N; i++ {
		FromResponse(401, "", false)
	}
}
