package web_scraper

import (
	"fmt"
	"log"
	"strings"
	"tgbotik/errs"
	"tgbotik/storage"
	"time"

	"github.com/gocolly/colly"
)

type Parser struct {
	c         *colly.Collector
	Domains   []string
	URLs      []string
	FirstTime bool
	Storage   *storage.PostgreDB
}

// Возвращает ссылку на объект типа Parser
func NewParser(domains, urls []string, firstTime bool, storage *storage.PostgreDB) *Parser {
	return &Parser{
		Domains:   domains,
		URLs:      urls,
		FirstTime: firstTime,
		Storage:   storage,
	}
}

// Запускает парсинг с переданными параметрами. В этом проекте данная функция
// должна запускаться параллельно
func Start(db *storage.PostgreDB, timeDelay time.Duration, URLs []string, domains []string) {
	parser := NewParser(domains, URLs, true, db)
	parser.CreateScraper()

	parser.Scrap()
	parser.FirstTime = false

	for t := range time.Tick(timeDelay) {
		log.Printf("Проверка наличия новых серий в %s...", t)
		parser.Scrap()
	}
}

// Создает colly collector как атрибут объекта типа Parser. Использует текущие параметры объекта типа Parser
func (p *Parser) CreateScraper() {
	c := colly.NewCollector(colly.AllowedDomains(p.Domains...))
	var animeName string
	var newSeriaHref string
	var newSeriaText string
	var countSeries, i, animeID int

	// Здесь извлекается название аниме
	c.OnHTML("h1.header_video", func(h *colly.HTMLElement) {
		animeName = strings.TrimPrefix(h.Text, "Смотреть ")
		animeName = strings.TrimSuffix(animeName, " все серии и сезоны")
	})
	var args []interface{}
	args = append(args, animeName)
	exist, err := p.Storage.Exists("Animes", []string{"AnimeName"}, args)
	errs.LogError(err)
	if exist {
		animeID, countSeries, err = p.Storage.CountSeries(animeName)
		errs.LogError(err)
	}

	// Здесь извлекается кол-во серий. В случае появления новой - извлекается информация о серии
	c.OnHTML("a.short-btn", func(element *colly.HTMLElement) {
		i += 1
		if i > countSeries {
			newSeriaHref = fmt.Sprintf("https://jut.su%s", element.Attr("href"))
			newSeriaText = fmt.Sprint(element.Text)
			countSeries = i
		}
	})

	if exist {
		err = p.Storage.SetNewSeries(animeID, newSeriaText, newSeriaHref)
		log.Printf("New seria in %s!", animeName)
	} else {
		err = p.Storage.SaveAnime(animeName, newSeriaText, newSeriaHref, countSeries)
		log.Printf("Finded new anime %s", animeName)
	}
	errs.LogError(err)

	c.OnRequest(collyOnRequest)
	c.OnError(collyOnError)

	p.c = c
}

// Запускает парсинг страниц
func (p *Parser) Scrap() {
	for _, url := range p.URLs {
		p.c.Visit(url)
	}
}

func collyOnError(response *colly.Response, err error) {
	fmt.Printf("Error while scraping: %s\n", err.Error())
}

func collyOnRequest(request *colly.Request) {
	request.Headers.Set("Accept-Language", "ru-RU;0.9")
	request.Headers.Set("Content-Encoding", "utf-8")
	fmt.Printf("Visiting %s\n", request.URL)
}
