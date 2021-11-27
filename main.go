package main

import (
	"fmt"
	"log"
	"strconv"
	"strings"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
)

func getKey() string {
	return "2118351815:AAGEdmU16piE7uD_7ojUMFZ5D1O4eQT1INk"
}

type wallet map[string]float64

var db = map[int64]wallet{}

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
		log.Println(update.Message.Text)

		msgArr := strings.Split(update.Message.Text, " ")

		switch msgArr[0] {
		case "ADD":
			summ, err := strconv.ParseFloat(msgArr[2], 64)
			if err != nil {
				bot.Send(tgbotapi.NewMessage(update.Message.Chat.ID, "Ошибка конвертации"))
				continue
			}

			if _, ok := db[update.Message.Chat.ID]; !ok {
				db[update.Message.Chat.ID] = wallet{}
			}

			db[update.Message.Chat.ID][msgArr[1]] += summ

			msg := fmt.Sprintf("Баланс: %s %f", msgArr[1], db[update.Message.Chat.ID][msgArr[1]])

			bot.Send(tgbotapi.NewMessage(update.Message.Chat.ID, msg))
		case "SUB":
			summ, err := strconv.ParseFloat(msgArr[2], 64)
			if err != nil {
				bot.Send(tgbotapi.NewMessage(update.Message.Chat.ID, "Ошибка конвертации"))
				continue
			}

			if _, ok := db[update.Message.Chat.ID]; !ok {
				db[update.Message.Chat.ID] = wallet{}
			}

			db[update.Message.Chat.ID][msgArr[1]] -= summ

			msg := fmt.Sprintf("Баланс: %s %f", msgArr[1], db[update.Message.Chat.ID][msgArr[1]])

			bot.Send(tgbotapi.NewMessage(update.Message.Chat.ID, msg))
		case "DEL":
			delete(db[update.Message.Chat.ID], msgArr[1])
			bot.Send(tgbotapi.NewMessage(update.Message.Chat.ID, "Валюта удалена"))
		case "SHOW":
			msg := "Баланс:\n"

			for key, value := range db[update.Message.Chat.ID] {
				msg += fmt.Sprintf("Валюта: %s Сумма: %f\n", key, value)
			}

			bot.Send(tgbotapi.NewMessage(update.Message.Chat.ID, msg))
		default:
			bot.Send(tgbotapi.NewMessage(update.Message.Chat.ID, "Команда неизвестна"))
		}

		// msg := tgbotapi.NewMessage(update.Message.Chat.ID, update.Message.Text)
		// msg.ReplyToMessageID = update.Message.MessageID

		// bot.Send(msg)
	}
}
