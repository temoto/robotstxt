What
====

This is a robots.txt exclusion protocol implementation for Go language (golang).


Build
=====

To build and run tests run `go test` in source directory.


Usage
=====

1. Parse
^^^^^^^^

First of all, you need to parse robots.txt data. You can do it with
function `FromString(body string) (*RobotsData, error)`::

    robots, err := robotstxt.FromString("User-agent: *\nDisallow:")

There is a convenient function `FromResponse(statusCode int, body string) (*RobotsData, error)`
to init robots data from HTTP response status code and body::

    robots, err := robotstxt.FromResponse(resp.StatusCode, resp.Body)
    if err != nil {
        // robots.txt parse error
        return false, err
    }

Passing status code applies following logic in line with google's interpretation of robots.txt files:

    * status code = 4xx      -> allow all (even 401/403, as recommended by Google).
    * status code = 2xx      -> parse body with `FromString` and apply rules listed there.
    * other statuses (5xx)   -> disallow all, consider this a temporary unavailability.

2. Query
^^^^^^^^

Parsing robots.txt content builds a kind of logic database, which you can
query with `(r *RobotsData) TestAgent(url, agent string) (bool)`.

Explicit passing of agent is useful if you want to query for different agents. For single agent
users there is a convenient option: `(r *RobotsData) Test(url) (bool)` which is
identical to `TestAgent`, but uses `r.DefaultAgent` as user agent for each query.

Query parsed robots data with explicit user agent.

::

    allow := robots.TestAgent("/", "FooBot")

Or with implicit user agent.

::

    robots.DefaultAgent = "OtherBot"
    allow := robots.Test("/")

