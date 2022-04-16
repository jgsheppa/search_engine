package scraper

import (
	"github.com/gocolly/colly/v2"
	"github.com/jgsheppa/search_engine/models"
	"strconv"
	"strings"
)

func ScrapeWebPage(url, htmlTag, containerClass string) models.Documents {
	var topic string
	var subHeaders []string
	var guides models.Documents
	baseUrl := "pkg.go.dev/"
	wholeUrl := baseUrl + url

	c := colly.NewCollector(
		colly.AllowedDomains("https://" + baseUrl),
	)
	c.OnHTML(containerClass, func(e *colly.HTMLElement) {
		topic = e.ChildText("h1")
		subHeaders = e.ChildTexts(htmlTag)
	})
	c.Visit(wholeUrl)

	for i, header := range subHeaders {
		headerWithoutEdit := strings.Split(header, "[")
		guides = append(guides, models.Document{
			Document: topic + strconv.Itoa(i) + "header",
			Text:     headerWithoutEdit[0],
			Topic:    topic,
			URL:      wholeUrl + "#" + headerWithoutEdit[0],
		})
	}

	return guides
}
