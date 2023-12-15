package fewdobby

import (
	"fmt"
	"log/slog"

	"github.com/gocolly/colly"
	"github.com/gocolly/colly/debug"
)

func Main() int {
	slog.Debug("fewdobby", "test", true)
	run()

	return 0
}

func run() {
	c := colly.NewCollector(
		colly.Debugger(&debug.LogDebugger{}),
		colly.AllowedDomains("gopro.com"),
		colly.UserAgent("Mozilla/5.0 (Windows NT 6.1) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/41.0.2228.0 Safari/537.36"),
	)

	// On every a element which has href attribute call callback
	c.OnHTML("a[href]", func(e *colly.HTMLElement) {
		link := e.Attr("href")
		// Print link
		fmt.Printf("Link found: %q -> %s\n", e.Text, link)
		// Visit link found on page
		// Only those links are visited which are in AllowedDomains
		c.Visit(e.Request.AbsoluteURL(link))
	})

	// Before making a request print "Visiting ..."
	c.OnRequest(func(r *colly.Request) {
		fmt.Println("Visiting", r.URL.String())
	})

	// start scripting here
	c.Visit("https://gopro.com/en/us/shop/cameras/hero12-black/CHDHX-121-master.html")
}
