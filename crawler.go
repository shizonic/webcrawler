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
	Documents []Document
}

func NewCrawler(startURL string, levels int) *Crawler {
	return &Crawler{
		WaitGroup: new(sync.WaitGroup),
		StartURL:  startURL,
		Levels:    levels,
		Crawled:   0,
		Documents: make([]Document, 0),
	}
}

func (c *Crawler) Crawle() {
	rootDocument := NewDocument(c.StartURL, 1)
	c.WaitGroup.Add(1)
	c.CrawleDocument(rootDocument)
	c.WaitGroup.Wait()
}

func (c *Crawler) ParseTitle(doc *goquery.Document, document *Document) {
	document.Title = doc.Find("title").Text()
}

func (c *Crawler) ParseH1(doc *goquery.Document, document *Document) {
	document.H1 = doc.Find("h1").Text()
}

func (c *Crawler) ParseBody(doc *goquery.Document, document *Document) {
	body, err := doc.Find("body").Html()
	if err != nil {
		logError(fmt.Sprintf(fmt.Sprintf("Failed (Lvl %v): %s", document.Level, document.URL)))
		return
	}
	document.Body = body
}

func (c *Crawler) ParseAnchors(doc *goquery.Document, document *Document) {
	doc.Find("a").Each(func(i int, anchor *goquery.Selection) {
		val, exists := anchor.Attr("href")
		if exists {
			document.URLs = append(document.URLs, val)

			c.WaitGroup.Add(1)
			go c.CrawleDocument(NewDocument(val, document.Level+1))
		}
	})
}

func (c *Crawler) CrawleDocument(document Document) {
	defer c.WaitGroup.Done()

	if document.Level <= c.Levels && strings.Index(document.URL, "http") == 0 {
		res, err := http.Get(document.URL)
		if err != nil {
			logError(fmt.Sprintf("Failed (Lvl %v): %s", document.Level, document.URL))
			return
		}
		defer res.Body.Close()

		doc, err := goquery.NewDocumentFromReader(res.Body)
		if err != nil {
			logError(fmt.Sprintf("Failed (Lvl %v): %s", document.Level, document.URL))
			return
		}

		c.ParseTitle(doc, &document)
		c.ParseH1(doc, &document)
		c.ParseBody(doc, &document)
		c.ParseAnchors(doc, &document)

		c.Documents = append(c.Documents, document)
		c.Crawled++

		logSuccess(fmt.Sprintf("Success (Lvl %v): %s", document.Level, document.URL))
	}
}
