package main

type Document struct {
	URL, Title, H1, Body string
	Level                int
	URLs                 []string
}

func NewDocument(url string, level int) Document {
	return Document{
		URL:   url,
		Level: level,
	}
}
