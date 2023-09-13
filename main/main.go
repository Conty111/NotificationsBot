package main

import (
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
	// Запускается ТГ бот, запоминает последние серии в БД (или их кол-во)
	// Запускается парсер (параллельно)
	// Если парсер находит новые серии, сохраняет их в БД и изменяет статус на новое
	// Тг бот проверяет статусы аниме, если новое -
	// извлекается последняя серия и animeID, меняется статус для аниме на проверено
	// По извлеченному animeID находим всех активных и подписанных на аниме и отправляем им серию

	// Система подписки: для подписки на аниме выбирает из перечня, после этого
	// заносится новая запись в Subscribers. Для отписки запись просто удаляется
	db, err := storage.New(host, user, password, dbname, port)
	errs.CheckError(err)
	defer func() {
		err = db.DB.Close()
		errs.CheckError(err)
	}()

	errs.CheckError(db.SaveUser(126, "dima", true))
	errs.CheckError(db.SaveUser(126, "dima", true))

	// var newSeries []string
	// go func() {
	// 	URLs := []string{"https://jut.su/mushoku-tensei/", "https://jut.su/kime-no-yaiba/"}
	// 	domains := []string{"www.jut.su", "jut.su"}
	// 	parser := web_scraper.NewParser(domains, URLs, true)
	// 	parser.CreateScraper()

	// 	timeDelay := time.Minute * 30

	// 	parser.Scrap()
	// 	parser.FirstTime = false

	// 	for now := range time.Tick(timeDelay) {
	// 		t := now.Format("2006-01-02 15:04")
	// 		fmt.Println(t, "Проверка наличия новых серий...")
	// 		Series := parser.Scrap()
	// 		if Series != nil {
	// 			fmt.Println(t, "Новая серия вышла!")
	// 			newSeries = Series
	// 		}
	// 	}
	// }()

	// bot, err := tg.New(os.Getenv("TGBOT_TOKEN"))
	// CheckError(err)
	// updateConfig := tgbotapi.NewUpdate(0)
	// updateConfig.Timeout = 15
	// updates := bot.GetUpdatesChan(updateConfig)

	// for update := range updates {
	// 	text := tg.CheckNewUpdate(update)
	// 	if text != "" {
	// 		log.Printf("Пользователь %s прислал сообщение: %s", update.Message.Chat.UserName, update.Message.Text)
	// 		msg := tgbotapi.NewMessage(update.Message.Chat.ID, text)
	// 		if _, err := bot.Send(msg); err != nil {
	// 			panic(err)
	// 		}
	// 	}

	// 	// Здесь проверять БД
	// 	if newSeries != nil {
	// 		for _, seria := range newSeries {
	// 			mes := fmt.Sprintf("Ееееее, новая серия вышла\n%s", seria)
	// 			log.Printf(mes, " отправлено пользователям")
	// 			msg := tgbotapi.NewMessage(update.Message.Chat.ID, mes)
	// 			if _, err := bot.Send(msg); err != nil {
	// 				panic(err)
	// 			}
	// 		}
	// 	}
	// }
}
