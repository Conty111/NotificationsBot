package tg

import (
	"fmt"
	"log"
	"strconv"
	"tgbotik/errs"
	"tgbotik/storage"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type Bot struct {
	TgClient *tgbotapi.BotAPI
	Storage  *storage.PostgreDB
}

// Returns the *Bot
func New(token string, s *storage.PostgreDB) (*Bot, error) {
	bot, err := tgbotapi.NewBotAPI(token)
	if err != nil {
		return nil, err
	}
	bot.Debug = false

	return &Bot{
		TgClient: bot,
		Storage:  s,
	}, nil
}

// Handle update from telegram
func (b *Bot) HandleUpdate(u tgbotapi.Update) error {
	if u.CallbackQuery != nil && u.CallbackQuery.Data != "" {
		return errs.CheckError(b.HandleCallback(u))
	}
	log.Printf("User %s sent the message: %s", u.Message.Text, u.Message.Chat.UserName)
	switch u.Message.Text {
	case start:
		err := b.Storage.SaveUser(int(u.Message.Chat.ID), u.Message.Chat.UserName, true)
		errs.LogError(err)
		return errs.CheckError(b.Response(printHello(u.Message.Chat.ID, u.Message.Chat.UserName)))
	case help:
		err := b.Response(printHelp(u.Message.Chat.ID))
		return errs.CheckError(err)
	case animes:
		err := b.Response(b.printAnimes(u.Message.Chat.ID))
		return errs.CheckError(err)
	default:
		err := b.Response(printDontSpam(u.Message.Chat.ID))
		return errs.CheckError(err)
	}
}

// Handle Callback data from telegram update
func (b *Bot) HandleCallback(u tgbotapi.Update) error {
	choiceAnimeID, err := strconv.Atoi(u.CallbackQuery.Data)
	errs.LogError(err)

	var args []interface{}
	args = append(args, u.CallbackQuery.Message.Chat.ID)
	args = append(args, choiceAnimeID)
	exist, err := b.Storage.Exists("Subscribers", []string{"ChatID", "AnimeID"}, args)
	var msg tgbotapi.MessageConfig
	if exist {
		b.Storage.Unsubscribe(int(u.CallbackQuery.Message.Chat.ID), choiceAnimeID)
		msg = tgbotapi.NewMessage(u.CallbackQuery.Message.Chat.ID, unsubcscribeText)
	} else {
		b.Storage.Subscribe(int(u.CallbackQuery.Message.Chat.ID), choiceAnimeID)
		msg = tgbotapi.NewMessage(u.CallbackQuery.Message.Chat.ID, subcscribeText)
	}
	return errs.CheckError(b.Response(msg))
}

// Send response to the chat
func (b *Bot) Response(msg tgbotapi.MessageConfig) error {
	log.Printf("Responsing '%s' to %d user", msg.Text, msg.ChatID)
	_, err := b.TgClient.Send(msg)
	return errs.CheckError(err)
}

// Returns message (type - tgbotapi.MessageConfig) that print a list of saved in DB anime
func (b *Bot) printAnimes(chatID int64) tgbotapi.MessageConfig {
	animeIDs, animeNames, err := b.Storage.Animes()
	errs.LogError(err)
	msg := tgbotapi.NewMessage(chatID, listAnimeTitle)
	var rows [][]tgbotapi.InlineKeyboardButton
	for idx, anime := range animeNames {
		rows = append(rows, getAnimeButtons(anime, animeIDs[idx]))
	}
	msg.ReplyMarkup = tgbotapi.NewInlineKeyboardMarkup(rows...)
	return msg
}

// Returns hello message
func printHello(chatID int64, userName string) tgbotapi.MessageConfig {
	return tgbotapi.NewMessage(
		chatID,
		helloText+fmt.Sprintf("%s!\nКакое из знакомых мне аниме ты смотришь?)",
			userName))
}

// Returns help message
func printHelp(chatID int64) tgbotapi.MessageConfig {
	return tgbotapi.NewMessage(chatID, helpText)
}

// Returns "don't spam" message
func printDontSpam(chatID int64) tgbotapi.MessageConfig {
	return tgbotapi.NewMessage(chatID, dontSpamText)
}

// Returns maked telegram inline button row for anime list
func getAnimeButtons(anime string, id int) []tgbotapi.InlineKeyboardButton {
	button := tgbotapi.NewInlineKeyboardButtonData(anime, fmt.Sprint(id))
	return tgbotapi.NewInlineKeyboardRow(button)
}
