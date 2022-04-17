package scraper

import (
	"github.com/gocolly/colly/v2"
	"github.com/jgsheppa/search_engine/models"
	"strconv"
)

func ScrapeWebPage(url, htmlTag, containerClass string) models.Documents {
	var topic string
	var subHeaders []string
	var guides models.Documents
	baseUrl := "pkg.go.dev"
	wholeUrl := "https://" + baseUrl + "/" + url

	c := colly.NewCollector(
		colly.AllowedDomains(baseUrl),
	)

	c.OnHTML(containerClass, func(e *colly.HTMLElement) {
		topic = e.ChildText("h1")
		subHeaders = e.ChildTexts(htmlTag)

	})
	c.Visit(wholeUrl)
	for i, header := range subHeaders {
		guides = append(guides, models.Document{
			Document: topic + ":header:" + strconv.Itoa(i),
			Text:     header,
			Topic:    topic,
			URL:      wholeUrl + "#" + header,
		})
	}

	return guides
}
