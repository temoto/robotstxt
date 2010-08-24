package robotstxt

import (
    "testing"
)


func TestFromResponseBasic(t *testing.T) {
    if _, err := FromResponse(200, ""); err != nil {
        t.Fatal("FromResponse MUST accept 200/\"\"")
    }
    if _, err := FromResponse(401, ""); err != nil {
        t.Fatal("FromResponse MUST accept 401/\"\"")
    }
    if _, err := FromResponse(403, ""); err != nil {
        t.Fatal("FromResponse MUST accept 403/\"\"")
    }
    if _, err := FromResponse(404, ""); err != nil {
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
    r, _ := FromResponse(401, "")
    ExpectDisallow(r, t, "FromResponse(401, \"\") MUST disallow everything.")
}

func TestResponse403(t *testing.T) {
    r, _ := FromResponse(403, "")
    ExpectDisallow(r, t, "FromResponse(403, \"\") MUST disallow everything.")
}

func TestResponse404(t *testing.T) {
    r, _ := FromResponse(404, "")
    ExpectAllow(r, t, "FromResponse(404, \"\") MUST allow everything.")
}


func TestFromStringBasic(t *testing.T) {
    if _, err := FromString(""); err != nil {
        t.Fatal("FromString MUST accept \"\"")
    }
}

func TestFromStringComment(t *testing.T) {
    if _, err := FromString("# comment"); err != nil {
        t.Fatal("FromString MUST accept \"# comment\"")
    }
}

func TestFromString001(t *testing.T) {
    r, err := FromString("User-Agent: *\r\nDisallow: /\r\n")
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
    r, err := FromString("User-Agent: *\r\nDisallow: /account\r\n")
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


func TestEmpty(t *testing.T) {
    r, _ := FromString("")
    if allow, err := r.TestAgent("/", "Somebot"); err != nil || !allow {
        t.Fatal("FromString(\"\") MUST allow everything.")
    }
}
