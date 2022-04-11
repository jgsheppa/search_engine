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
	var header4 []string
	var guides models.Guides

	c := colly.NewCollector(
		colly.AllowedDomains("guidebook.bestpracticer.com"),
	)
	c.OnHTML(".content", func(e *colly.HTMLElement) {
		topic = e.ChildText("h2")
		subHeaders = e.ChildTexts("h3")
		header4 = e.ChildTexts("h4")
	})
	c.Visit(baseURL)

	for i, header := range subHeaders {
		guides = append(guides, models.Guide{
			Document: topic + strconv.Itoa(i) + "header",
			Text:     header,
			Topic:    topic,
			URL:      baseURL,
		})
	}

	for i, header := range header4 {
		guides = append(guides, models.Guide{
			Document: topic + strconv.Itoa(i) + "header4",
			Text:     header,
			Topic:    topic,
			URL:      baseURL,
		})
	}
	return guides
}
