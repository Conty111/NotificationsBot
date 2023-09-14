package main

import (
	"fmt"
	"log"
	"os"
	"tgbotik/errs"
	"tgbotik/storage"
	"tgbotik/tg"
	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

const (
	host           = "localhost"
	port           = 5432
	user           = "bot"
	password       = "Mushoku777!"
	dbname         = "Rudeus"
	timeParseDelay = time.Minute * 30
)

var domains = []string{"jut.su", "www.jut.su"}

func main() {
	// Подключаемся к БД
	// Запускается парсер в новом потоке
	// Запускается ТГ бот в основном потоке
	// Если парсер находит новые серии, сохраняет их в БД и изменяет статус для аниме новое (true)
	// Тг бот вместе с получениями update-ов проверяет аниме в БД, если находит со статусом true, то
	// извлекается последняя серия и animeID, меняется статус для аниме на false
	// По извлеченному animeID находим всех активных и подписанных на аниме пользователей
	// Проходясь по полученному списку пользователей, отправляем им сообщения в ТГ

	db, err := storage.New(host, user, password, dbname, port)
	errs.LogError(err)
	defer func() {
		err = db.DB.Close()
		errs.LogError(err)
	}()

	db.SaveAnime("TestAnime", "Last Seria", "href", 56)
	db.SaveAnime("TestAnime2", "Last Seria2", "href2", 57)

	// Запускается парсер
	// go web_scraper.Start(db, timeParseDelay, []string{}, domains)

	bot, err := tg.New(os.Getenv("TGBOT_TOKEN"), db)
	errs.LogError(err)

	updateConfig := tgbotapi.NewUpdate(0)
	updateConfig.Timeout = 15
	updates := bot.TgClient.GetUpdatesChan(updateConfig)

	log.Print("Started get updates from telegram API")
	// Запускается прослушивание входящих update-ов
	for update := range updates {
		bot.HandleUpdate(update)

		// Здесь проверяется БД на наличие новых серий
		newSeries, err := db.GetNewSeries()
		errs.LogError(err)
		if len(newSeries) > 0 {
			log.Print("Finded new series")
			for _, s := range newSeries {
				text := fmt.Sprintf("Хэй-хэй-хэй, вышла новая серия по аниме %s\nА вот и она: %s\n%s",
					s.AnimeName, s.Text, s.Href)

				// Получаем пользователей, подписанных на аниме и отправляем им сообщения
				usersID, err := db.GetSubscribers(s.ID)
				errs.LogError(err)
				for _, ID := range usersID {
					msg := tgbotapi.NewMessage(int64(ID), text)
					_, err = bot.TgClient.Send(msg)
					errs.LogError(err)
				}
				db.SetStatus("Animes", s.ID, false)
			}
		}
	}
}
