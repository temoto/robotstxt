package main

import (
	"flag"
	"log"
	"net/http"
	"strings"

	"github.com/temoto/robotstxt"
)

var checkPaths = []string{
	"/",
}

func main() {
	robotsUrl := flag.String("robots-url", "", "")
	bot := flag.String("bot", "GoogleBot", "")
	flag.Parse()
	if *robotsUrl == "" {
		log.Fatalln("Robots URL is empty, run with -h to see usage.")
	}
	if !strings.HasPrefix(*robotsUrl, "http") {
		*robotsUrl = "http://" + *robotsUrl
	}

	response, err := http.Get(*robotsUrl)
	if err != nil {
		log.Fatalln("HTTP error:", err)
	}

	robots, err := robotstxt.FromResponse(response)
	if err != nil {
		log.Fatalln("Robots.txt error:", err)
	}

	log.Println("Running checks as", *bot)
	group := robots.FindGroup(*bot)
	for _, path := range checkPaths {
		log.Println(path, ":", group.Test(path))
	}
}
