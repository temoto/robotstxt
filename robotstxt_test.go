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
    allow, err := r.TestAgent("/", "Somebot")
    if err != nil {
        t.Fatal("Unexpected error.")
    }
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
    ExpectDisallow(r, t, "FromResponse(401, \"\") MUST disallow everything.")
}

func TestResponse403(t *testing.T) {
    r, _ := FromResponse(403, "", true)
    ExpectDisallow(r, t, "FromResponse(403, \"\") MUST disallow everything.")
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
    if allow, err := r.TestAgent("/", "Somebot"); err != nil || !allow {
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
        t.Fatal(err.String())
    }
    allow, err1 := r.TestAgent("/foobar", "SomeAgent")
    if err1 != nil {
        t.Fatal(err1.String())
    }
    if allow {
        t.Fatal("Must deny.")
    }
}

func TestFromString002(t *testing.T) {
    r, err := FromString("User-Agent: *\r\nDisallow: /account\r\n", true)
    if err != nil {
        t.Fatal(err.String())
    }
    allow, err1 := r.TestAgent("/foobar", "SomeAgent")
    if err1 != nil {
        t.Fatal(err1.String())
    }
    if !allow {
        t.Fatal("Must allow.")
    }
}

const robots_text_001 = "User-agent: * \nDisallow: /administrator/\nDisallow: /cache/\nDisallow: /components/\nDisallow: /editor/\nDisallow: /forum/\nDisallow: /help/\nDisallow: /images/\nDisallow: /includes/\nDisallow: /language/\nDisallow: /mambots/\nDisallow: /media/\nDisallow: /modules/\nDisallow: /templates/\nDisallow: /installation/\nDisallow: /getcid/\nDisallow: /tooltip/\nDisallow: /getuser/\nDisallow: /download/\nDisallow: /index.php?option=com_phorum*,quote=1\nDisallow: /index.php?option=com_phorum*phorum_query=search\nDisallow: /index.php?option=com_phorum*,newer\nDisallow: /index.php?option=com_phorum*,older\n\nUser-agent: Yandex\nAllow: /\nSitemap: http://www.pravorulya.com/sitemap.xml\nSitemap: http://www.pravorulya.com/sitemap1.xml"

func TestFromString003(t *testing.T) {
    r, err := FromString(robots_text_001, true)
    if err != nil {
        t.Fatal(err.String())
    }
    allow, err1 := r.TestAgent("/administrator/", "SomeBot")
    if err1 != nil {
        t.Fatal(err1.String())
    }
    if allow {
        t.Fatal("Must deny.")
    }
}

func TestFromString004(t *testing.T) {
    r, err := FromString(robots_text_001, true)
    if err != nil {
        t.Fatal(err.String())
    }
    allow, err1 := r.TestAgent("/paruram", "SomeBot")
    if err1 != nil {
        t.Fatal(err1.String())
    }
    if !allow {
        t.Fatal("Must allow.")
    }
}
