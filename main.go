package main

import (
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/gocolly/colly"
)

func main() {
	total := 0
	count := 0

	url := os.Args[1]
	if len(url) < 1 {
		log.Fatal("Please enter a URL to crawl.")
	}

	parts := strings.Split(url, "//")
	if len(parts) < 2 {
		log.Fatal("Please enter a valid URL (including scheme)")
	}
	domain := parts[1]

	c := colly.NewCollector(
		colly.AllowedDomains(domain),
	)

	c.OnHTML("a[href]", func(e *colly.HTMLElement) {
		e.Request.Visit(e.Attr("href"))
	})

	c.OnRequest(func(r *colly.Request) {
		r.Ctx.Put("start", fmt.Sprint(time.Now().UnixNano()))
	})

	c.OnResponse(func(r *colly.Response) {
		s := r.Ctx.Get("start")
		start, err := strconv.Atoi(s)
		if err != nil {
			log.Fatal("start is not an int")
		}

		now := int(time.Now().UnixNano())

		elapsed := (now - start) / 1000000
		total += elapsed
		count++
		average := total / count

		fmt.Printf("%d\t%s\t%d ms\tavg: %d\n", r.StatusCode, r.Request.URL, elapsed, average)
	})

	c.Visit(url)
}
