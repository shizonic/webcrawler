package main

import (
	"fmt"
	"os"
	"strconv"
	"sync"
)

var wg sync.WaitGroup

func main() {
	url := os.Args[1]
	levels, err := strconv.Atoi(os.Args[2])
	if err != nil {
		panic(err)
	}

	crawler := NewCrawler(&wg, url, levels)
	wg.Add(1)
	fmt.Printf("===\nStart crawling from: %s\n===\n", url)
	go crawler.Crawle()
	wg.Wait()
	fmt.Printf("===\n%v crawled\n===\n", crawler.Crawled)
}
