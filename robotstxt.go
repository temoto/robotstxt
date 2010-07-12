// The robots.txt Exclusion Protocol is implemented as specified in
// http://www.robotstxt.org/wc/robots.html
// with various extensions.
package robotstxt

import (
    "os"
    "strings"
)


type RobotsData struct {
    DefaultAgent string
    // private
    data         string
    allowAll     bool
    disallowAll  bool
}

var AllowAll = &RobotsData{allowAll: true}
var DisallowAll = &RobotsData{disallowAll: true}


func FromResponse(statusCode int, body string) (*RobotsData, os.Error) {
    switch {
    case statusCode == 404:
        return AllowAll, nil
    case statusCode == 401 || statusCode == 403:
        return DisallowAll, nil
    case statusCode >= 200 && statusCode < 300:
        return FromString(body)
    }
    // Conservative disallow all default
    return DisallowAll, nil
}

func FromString(body string) (*RobotsData, os.Error) {
    trimmed := strings.TrimSpace(body)
    if trimmed == "" {
        return AllowAll, nil
    }

    return nil, os.NewError("TODO: parsing non-empty robots.txt is not implemented yet")
    return DisallowAll, nil
}

func (r *RobotsData) Test(url string) (bool, os.Error) {
    if r.DefaultAgent == "" {
        return false, os.NewError("DefaultAgent is empty. You MUST set RobotsData.DefaultAgent to use Test method.")
    }
    return r.TestAgent(url, r.DefaultAgent)
}

func (r *RobotsData) TestAgent(url, agent string) (bool, os.Error) {
    if r.allowAll {
        return true, nil
    }
    if r.disallowAll {
        return false, nil
    }
    return false, nil
}
