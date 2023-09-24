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

type Anime struct {
	ID            int
	Name          string
	LastSeriaText string
	LastSeriaHref string
	CountSeries   int
	Status        bool
}

var a *Anime
var series [][]string

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
	log.Print("Started to work the parser")
	parser := NewParser(domains, URLs, true, db)
	parser.CreateScraper()

	parser.Scrap()

	for t := range time.Tick(timeDelay) {
		log.Printf("Проверка наличия новых серий в %s...", t)
		parser.Scrap()
	}
}

// Создает colly collector как атрибут объекта типа Parser. Использует текущие параметры объекта типа Parser
func (p *Parser) CreateScraper() {
	c := colly.NewCollector(colly.AllowedDomains(p.Domains...))
	var i int

	// Здесь извлекается название аниме
	c.OnHTML("h1.header_video", func(h *colly.HTMLElement) {
		log.Print("Fetching the anime name")
		animeName := strings.TrimPrefix(h.Text, "Смотреть ")
		animeName = strings.TrimSuffix(animeName, " все серии и сезоны")
		exist, err := p.CheckAnime(animeName)
		errs.LogError(err)
		if exist {
			animeID, countSeries, err := p.Storage.CountSeries(animeName)
			errs.LogError(err)
			a.CountSeries = countSeries
			a.ID = animeID
		}
		a.Name = animeName
	})

	// Здесь извлекается кол-во серий. В случае появления новой - извлекается информация о серии
	c.OnHTML("a.short-btn", func(element *colly.HTMLElement) {
		i += 1
		if i > a.CountSeries {
			a.LastSeriaHref = fmt.Sprintf("https://jut.su%s", element.Attr("href"))
			a.LastSeriaText = fmt.Sprint(element.Text)
			a.CountSeries = i
			a.Status = true
		}
	})

	c.OnRequest(collyOnRequest)
	c.OnError(collyOnError)

	p.c = c
}

// Запускает парсинг страниц
func (p *Parser) Scrap() {
	for _, url := range p.URLs {
		p.c.Visit(url)
		if a.ID == 0 {
			errs.LogError(p.Storage.SaveAnime(a.Name, a.LastSeriaText, a.LastSeriaHref, a.CountSeries))
		} else if a.Status {
			errs.LogError(p.Storage.SetNewSeries(a.ID, a.LastSeriaText, a.LastSeriaHref))
		}
	}
}

// Check if exist the anime by name
func (p *Parser) CheckAnime(animeName string) (bool, error) {
	var args []interface{}
	args = append(args, animeName)
	log.Print("Checking saved anime")
	exist, err := p.Storage.Exists("Animes", []string{"AnimeName"}, args)
	return exist, err
}

func collyOnError(response *colly.Response, err error) {
	fmt.Printf("Error while scraping: %s\n", err.Error())
}

func collyOnRequest(request *colly.Request) {
	request.Headers.Set("Accept-Language", "ru-RU;0.9")
	request.Headers.Set("Content-Encoding", "utf-8")
	fmt.Printf("Visiting %s\n", request.URL)
}
