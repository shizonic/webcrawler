package main

type Site struct {
	URL, Title, H1, Body string
	Level                int
	URLs                 []string
}

func NewSite(url string, level int) Site {
	return Site{
		URL:   url,
		Level: level,
	}
}
