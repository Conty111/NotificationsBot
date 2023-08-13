package web_scraper

import (
	"fmt"
	"github.com/gocolly/colly"
	"strings"
)

var animeSeriesCount = make(map[string]int)
var newSerias []string

type Parser struct {
	c         *colly.Collector
	Domains   []string
	URLs      []string
	FirstTime bool
}

func NewParser(domains, urls []string, firstTime bool) *Parser {
	return &Parser{
		Domains:   domains,
		URLs:      urls,
		FirstTime: firstTime,
	}
}

func (p *Parser) CreateScraper() {
	c := colly.NewCollector(colly.AllowedDomains(p.Domains...))
	var countSeries int
	var animeName string

	c.OnHTML("h1.header_video", func(h *colly.HTMLElement) {
		animeName = strings.TrimPrefix(h.Text, "Смотреть ")
		animeName = strings.TrimSuffix(animeName, " все серии и сезоны")
	})
	c.OnHTML("a.short-btn", func(element *colly.HTMLElement) {
		countSeries += 1
		if countSeries > animeSeriesCount[animeName] && !p.FirstTime {
			msg := fmt.Sprintf("%s: https://jut.su%s", element.Text, element.Attr("href"))
			fmt.Println(msg)
			newSerias = append(newSerias, msg)
			animeSeriesCount[animeName] = countSeries
		} else if p.FirstTime {
			animeSeriesCount[animeName] = countSeries
		}
	})
	c.OnRequest(collyOnRequest)
	c.OnError(collyOnError)
	p.c = c
}

func (p *Parser) Scrap() []string {
	for _, url := range p.URLs {
		p.c.Visit(url)
	}
	return newSerias
}

func collyOnError(response *colly.Response, err error) {
	fmt.Printf("Error while scraping: %s\n", err.Error())
}

func collyOnRequest(request *colly.Request) {
	request.Headers.Set("Accept-Language", "ru-RU;0.9")
	request.Headers.Set("Content-Encoding", "utf-8")
	fmt.Printf("Visiting %s\n", request.URL)
}
