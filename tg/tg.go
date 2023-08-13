package tg

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func New(token string) (*tgbotapi.BotAPI, error) {
	bot, err := tgbotapi.NewBotAPI(token)
	if err != nil {
		return nil, err
	}
	bot.Debug = true
	return bot, nil
}

func CheckNewUpdate(update tgbotapi.Update) string {
	if update.Message.Text == "/start" {
		return "Вы подписались на рассылку"
	} else if update.Message.Text != "" {
		return "А вот спамить не надо"
	}
	return ""
}
