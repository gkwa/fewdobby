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
	collector := colly.NewCollector(
		colly.Debugger(&debug.LogDebugger{}),
		colly.AllowedDomains("gopro.com"),
		colly.UserAgent("Mozilla/5.0 (Windows NT 6.1) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/41.0.2228.0 Safari/537.36"),
	)
	collector.IgnoreRobotsTxt = true

	// On every a element which has href attribute call callback
	collector.OnHTML("a[href]", func(e *colly.HTMLElement) {
		link := e.Attr("href")
		fmt.Printf("Link found: %q -> %s\n", e.Text, link)
		collector.Visit(e.Request.AbsoluteURL(link))
	})

	// On error, print the error and response if available
	collector.OnError(func(r *colly.Response, err error) {
		if r != nil && r.Request != nil {
			fmt.Printf("Error visiting %s: %v\n", r.Request.URL.String(), err)
			fmt.Printf("Response from %s:\n%s\n", r.Request.URL.String(), r.Body)
		} else if r != nil {
			fmt.Printf("Error during request: %v\n", err)
			fmt.Printf("Response from an unknown URL:\n%s\n", r.Body)
		} else {
			fmt.Printf("Error during request: %v\n", err)
		}
	})

	// Before making a request print "Visiting ..."
	collector.OnRequest(func(r *colly.Request) {
		fmt.Println("Visiting", r.URL.String())
	})

	collector.Limit(&colly.LimitRule{DomainGlob: "*", Parallelism: 1})

	// On response, print the response body as text
	collector.OnResponse(func(r *colly.Response) {
		fmt.Printf("Response from %s:\n%s\n", r.Request.URL.String(), r.Body)
	})

	// start scripting here
	collector.Visit("https://gopro.com/en/us/shop/cameras/hero12-black/CHDHX-121-master.html")

	// Wait until threads are finished
	collector.Wait()
}
