package tg

import tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"

const (
	sparkles_em  = "\U00002728"
	heart_em     = "\U00002764"
	arrDown_em   = "\U00002b07"
	x_em         = "\U0000274c"
	checkMark_em = "\U00002705"
	cry_em       = "\U0001F972"
	class_em     = "\U0001F44D"

	start            = "/start"
	help             = "Help me " + cry_em
	animes           = heart_em + " Список аниме " + heart_em
	listAnimeTitle   = arrDown_em + " Список доступных аниме ниже " + arrDown_em
	helloText        = "Пипипупу" + sparkles_em + ", приветствую "
	helpText         = "Помоги себе сам, друг)))\nНет, ну что тут непонятного то - есть список доступных мне аниме. Нажимаешь на любое(-ые) - и подписываешься на оповещения о выходе новых серий " + class_em
	subcscribeText   = checkMark_em + " Вы успешно подписались на оповещения " + checkMark_em
	unsubcscribeText = "Вы успешно отписались от оповещений об аниме. "
	dontSpamText     = x_em + " НЕ СПАМИТЬ, ИСПОЛЬЗОВАТЬ КНОПКИ!!! " + x_em
)

var HELP tgbotapi.KeyboardButton = tgbotapi.NewKeyboardButton(help)
var ANIME tgbotapi.KeyboardButton = tgbotapi.NewKeyboardButton(animes)

var row1 []tgbotapi.KeyboardButton = tgbotapi.NewKeyboardButtonRow(ANIME, HELP)
var keyboard tgbotapi.ReplyKeyboardMarkup = tgbotapi.NewReplyKeyboard(row1)
