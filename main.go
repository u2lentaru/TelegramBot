package main

import (
	"log"
	"strings"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
)

func getKey() string {
	return "2118351815:AAGEdmU16piE7uD_7ojUMFZ5D1O4eQT1INk"
}

func main() {
	//TestTB TestTB2bot
	bot, err := tgbotapi.NewBotAPI(getKey())
	if err != nil {
		log.Panic(err)
	}

	// bot.Debug = true

	log.Printf("Authorized on account %s", bot.Self.UserName)

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates, err := bot.GetUpdatesChan(u)

	for update := range updates {
		if update.Message == nil { // ignore any non-Message Updates
			continue
		}

		// log.Printf("[%s] %s", update.Message.From.UserName, update.Message.Text)

		msgArr := strings.Split(update.Message.Text, " ")

		switch msgArr[0] {
		case "ADD":
			bot.Send(tgbotapi.NewMessage(update.Message.Chat.ID, "Добавить валюту"))
		case "SUB":
			bot.Send(tgbotapi.NewMessage(update.Message.Chat.ID, "Убавить валюту"))
		case "DEL":
			bot.Send(tgbotapi.NewMessage(update.Message.Chat.ID, "Удалить валюту"))
		case "SHOW":
			bot.Send(tgbotapi.NewMessage(update.Message.Chat.ID, "Показать всё"))
		default:
			bot.Send(tgbotapi.NewMessage(update.Message.Chat.ID, "Команда неизвестна"))
		}

		// msg := tgbotapi.NewMessage(update.Message.Chat.ID, update.Message.Text)
		// msg.ReplyToMessageID = update.Message.MessageID

		// bot.Send(msg)
	}
}
