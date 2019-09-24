package main

import (
	"fmt"
	"net/http"
	"strings"
	"sync"

	"github.com/PuerkitoBio/goquery"
)

type Crawler struct {
	WaitGroup *sync.WaitGroup
	StartURL  string
	Levels    int
	Crawled   int
	Sites     []Site
}

func NewCrawler(startURL string, levels int) *Crawler {
	return &Crawler{
		WaitGroup: new(sync.WaitGroup),
		StartURL:  startURL,
		Levels:    levels,
		Crawled:   0,
		Sites:     make([]Site, 0),
	}
}

func (c *Crawler) Crawle() {
	rootSite := NewSite(c.StartURL, 1)
	c.WaitGroup.Add(1)
	c.CrawleSite(rootSite)
	c.WaitGroup.Wait()
}

func (c *Crawler) ParseTitle(doc *goquery.Document, site *Site) {
	site.Title = doc.Find("title").Text()
}

func (c *Crawler) ParseH1(doc *goquery.Document, site *Site) {
	site.H1 = doc.Find("h1").Text()
}

func (c *Crawler) ParseBody(doc *goquery.Document, site *Site) {
	body, err := doc.Find("body").Html()
	if err != nil {
		logError(fmt.Sprintf(fmt.Sprintf("Failed (Lvl %v): %s", site.Level, site.URL)))
		return
	}
	site.Body = body
}

func (c *Crawler) ParseAnchors(doc *goquery.Document, site *Site) {
	doc.Find("a").Each(func(i int, anchor *goquery.Selection) {
		val, exists := anchor.Attr("href")
		if exists {
			site.URLs = append(site.URLs, val)

			c.WaitGroup.Add(1)
			go c.CrawleSite(NewSite(val, site.Level+1))
		}
	})
}

func (c *Crawler) CrawleSite(site Site) {
	defer c.WaitGroup.Done()

	if site.Level <= c.Levels && strings.Index(site.URL, "http") == 0 {
		res, err := http.Get(site.URL)
		if err != nil {
			logError(fmt.Sprintf("Failed (Lvl %v): %s", site.Level, site.URL))
			return
		}
		defer res.Body.Close()

		doc, err := goquery.NewDocumentFromReader(res.Body)
		if err != nil {
			logError(fmt.Sprintf("Failed (Lvl %v): %s", site.Level, site.URL))
			return
		}

		c.ParseTitle(doc, &site)
		c.ParseH1(doc, &site)
		c.ParseBody(doc, &site)
		c.ParseAnchors(doc, &site)

		c.Sites = append(c.Sites, site)
		c.Crawled++

		logSuccess(fmt.Sprintf("Success (Lvl %v): %s", site.Level, site.URL))
	}
}
