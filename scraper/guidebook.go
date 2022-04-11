package scraper

import (
	"github.com/gocolly/colly/v2"
	"github.com/jgsheppa/search_engine/models"
	"strconv"
)

func ScrapeGuidebookLinks() []string {
	baseURL := "https://guidebook.bestpracticer.com/"
	var urlPaths []string

	c := colly.NewCollector(
		colly.AllowedDomains("guidebook.bestpracticer.com"),
	)
	c.OnHTML(".div-block-18-copy-copy", func(e *colly.HTMLElement) {
		links := e.ChildAttrs("a", "href")
		for _, link := range links {
			urlPaths = append(urlPaths, baseURL+link)
		}
	})

	c.Visit(baseURL)

	return urlPaths
}

func ScrapeGuidebookPages(url string) models.Guides {
	baseURL := "https://guidebook.bestpracticer.com/" + url
	var topic string
	var subHeaders []string
	var paragraphs []string
	var guides models.Guides

	c := colly.NewCollector(
		colly.AllowedDomains("guidebook.bestpracticer.com"),
	)
	c.OnHTML(".content", func(e *colly.HTMLElement) {
		topic = e.ChildText("h2")
		subHeaders = e.ChildTexts("h3")
		paragraphs = e.ChildTexts("div")
	})
	c.Visit(baseURL)

	for i, header := range subHeaders {
		guides = append(guides, models.Guide{
			Document: topic + strconv.Itoa(i) + "header",
			Text:     header,
			Topic:    topic,
			URL:      url,
		})
	}

	for i, paragraph := range paragraphs {
		guides = append(guides, models.Guide{
			Document: topic + strconv.Itoa(i) + "paragraph",
			Text:     paragraph,
			Topic:    topic,
			URL:      url,
		})
	}
	return guides
}
