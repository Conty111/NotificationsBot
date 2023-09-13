package main

import (
	"fmt"
	"log"
	"tgbotik/errs"
	"tgbotik/storage"
)

const (
	host     = "localhost"
	port     = 5432
	user     = "bot"
	password = "Mushoku777!"
	dbname   = "Rudeus"
)

func main() {
	// Подключаемся к БД
	// Запускается ТГ бот
	// Запускается парсер (параллельно)
	// Если парсер находит новые серии, сохраняет их в БД и изменяет статус на новое
	// Тг бот проверяет статусы аниме, если новое -
	// извлекается последняя серия и animeID, меняется статус для аниме на проверено
	// По извлеченному animeID находим всех активных и подписанных на аниме и отправляем им серию
	// Система подписки: для подписки на аниме выбирает из перечня, после этого
	// заносится новая запись в Subscribers. Для отписки запись просто удаляется

	db, err := storage.New(host, user, password, dbname, port)
	errs.LogError(err)
	defer func() {
		err = db.DB.Close()
		errs.LogError(err)
	}()
	errs.LogError(db.SaveUser(111, "SomeUser", true))
	log.Print(db.SaveAnime("SomeAnimer", "text", "href", 6))
	errs.LogError(db.Subscribe(111, 15))
	errs.LogError(db.Unsubscribe(111, 15))
	errs.LogError(db.Subscribe(111, 15))
	errs.LogError(db.Subscribe(111, 15))
	fmt.Println(db.CountSeries("SomeAnimer"))

	// go func() {
	// 	URLs := []string{"https://jut.su/mushoku-tensei/", "https://jut.su/kime-no-yaiba/"}
	// 	domains := []string{"www.jut.su", "jut.su"}
	// 	parser := web_scraper.NewParser(domains, URLs, true, db)
	// 	parser.CreateScraper()

	// 	timeDelay := time.Minute * 30

	// 	parser.Scrap()
	// 	parser.FirstTime = false

	// 	for _ = range time.Tick(timeDelay) {
	// 		log.Print("Проверка наличия новых серий...")
	// 		parser.Scrap()
	// 	}
	// }()

	// bot, err := tg.New(os.Getenv("TGBOT_TOKEN"))
	// errs.LogError(err)
	// updateConfig := tgbotapi.NewUpdate(0)
	// updateConfig.Timeout = 15
	// updates := bot.GetUpdatesChan(updateConfig)

	// for update := range updates {
	// 	text := tg.CheckNewUpdate(update)
	// 	if text != "" {
	// 		log.Printf("Пользователь %s прислал сообщение: %s", update.Message.Chat.UserName, update.Message.Text)
	// 		msg := tgbotapi.NewMessage(update.Message.Chat.ID, text)
	// 		_, err = bot.Send(msg)
	// 		errs.LogError(err)
	// 	}

	// 	// Здесь проверять БД
	// 	new, err := db.GetNewSeries()
	// 	errs.LogError(err)
	// 	if len(new) > 0 {
	// 		for _, elem := range new {
	// 			text := fmt.Sprintf("Хэй-хэй-хэй, вышла новая серия по аниме %s\nА вот и она: %s %s",
	// 				elem.AnimeName, elem.Text, elem.Href)
	// 			usersID, err := db.GetSubscribers(elem.ID)
	// 			errs.LogError(err)
	// 			for _, ID := range usersID {
	// 				msg := tgbotapi.NewMessage(ID, text)
	// 				_, err = bot.Send(msg)
	// 				errs.LogError(err)
	// 			}
	// 			db.SetStatus("Animes", elem.ID, false)
	// 		}
	// 	}
	// }
}
