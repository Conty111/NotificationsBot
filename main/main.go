package main

import (
	"math/rand"
	"sync"
	"tgbotik/errs"
	"tgbotik/storage"
	"tgbotik/web_scraper"
	"time"
)

const (
	host           = "localhost"
	port           = 5432
	user           = "bot"
	password       = "Mushoku777!"
	dbname         = "Rudeus"
	timeParseDelay = time.Minute * 15
)

var domains = []string{"jut.su", "www.jut.su"}
var sites = []string{"https://jut.su/mushoku-tensei/",
	"https://jut.su/shingekii-no-kyojin/",
	"https://jut.su/kime-no-yaiba/",
	"https://jut.su/full-metal-alchemist/"}

func main() {
	// Подключаемся к БД
	// Запускается парсер в новом потоке
	// Запускается ТГ бот в основном потоке
	// Если парсер находит новые серии, сохраняет их в БД и изменяет статус для аниме новое (true)
	// Тг бот вместе с получениями update-ов проверяет аниме в БД, если находит со статусом true, то
	// извлекается последняя серия и animeID, меняется статус для аниме на false
	// По извлеченному animeID находим всех активных и подписанных на аниме пользователей
	// Проходясь по полученному списку пользователей, отправляем им сообщения в ТГ

	rand.Seed(time.Now().UnixNano()) // Удалить потом

	db, err := storage.New(host, user, password, dbname, port)
	errs.LogError(err)
	defer func() {
		err = db.DB.Close()
		errs.LogError(err)
	}()

	var wg sync.WaitGroup // Удалить потом

	// Запускается парсер
	wg.Add(1) // Удалить потом
	go web_scraper.Start(db, timeParseDelay, sites, domains)
	wg.Wait() // Удалить потом

	// bot, err := tg.New(os.Getenv("TGBOT_TOKEN"), db)
	// errs.LogError(err)

	// updateConfig := tgbotapi.NewUpdate(0)
	// updateConfig.Timeout = 15
	// updates := bot.TgClient.GetUpdatesChan(updateConfig)

	// log.Print("Started get updates from telegram API")
	// // Запускается прослушивание входящих update-ов
	// for update := range updates {
	// 	bot.HandleUpdate(update)

	// 	// Здесь проверяется БД на наличие новых серий
	// 	newSeries, err := db.GetNewSeries()
	// 	errs.LogError(err)
	// 	if len(newSeries) > 0 {
	// 		log.Print("Finded new series")
	// 		for _, s := range newSeries {
	// 			text := fmt.Sprintf("Хэй-хэй-хэй, вышла новая серия по аниме %s\nА вот и она: %s\n%s",
	// 				s.AnimeName, s.Text, s.Href)

	// 			// Получаем пользователей, подписанных на аниме и отправляем им сообщения
	// 			usersID, err := db.GetSubscribers(s.ID)
	// 			errs.LogError(err)
	// 			for _, ID := range usersID {
	// 				msg := tgbotapi.NewMessage(int64(ID), text)
	// 				_, err = bot.TgClient.Send(msg)
	// 				errs.LogError(err)
	// 			}
	// 			db.SetStatus("Animes", s.ID, false)
	// 		}
	// 	}
	// }
}
