package main

import (
	"fmt"
	"os"
	"strconv"
)

func logError(msg string) {
	fmt.Printf("!> %s.\n", msg)
}

func logSuccess(msg string) {
	fmt.Printf("=> %s.\n", msg)
}

func die(msg string) {
	logError(msg)
	os.Exit(0)
}

func main() {
	if len(os.Args) < 3 {
		die("At least 2 arguments required")
	}
	url := os.Args[1]
	levels, err := strconv.Atoi(os.Args[2])
	if err != nil {
		die("Second argument has to be an integer")
	}

	crawler := NewCrawler(url, levels)
	fmt.Printf("===\nStart crawling from: %s\n===\n", url)
	crawler.Crawle()
	fmt.Printf("===\n%v crawled\n===\n", crawler.Crawled)
}
