package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
)

type binanceResp struct {
	Price float64 `json:"price,string"`
	Code  int64   `json:"code"`
}

func getKey() string {
	return "2118351815:AAGEdmU16piE7uD_7ojUMFZ5D1O4eQT1INk"
}

func getPrice(coin string) (price float64, err error) {
	resp, err := http.Get(fmt.Sprintf("https://api.binance.com/api/v3/ticker/price?symbol=%sUSDT", coin))

	if err != nil {
		return
	}
	defer resp.Body.Close()

	var jsonResp binanceResp

	err = json.NewDecoder(resp.Body).Decode(&jsonResp)
	if err != nil {
		return
	}

	if jsonResp.Code != 0 {
		err = errors.New("Некорректная валюта")
		return
	}

	price = jsonResp.Price

	return
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
			var usdSum float64 = 0

			for key, value := range db[update.Message.Chat.ID] {
				coinPrice, err := getPrice(key)

				if err != nil {
					bot.Send(tgbotapi.NewMessage(update.Message.Chat.ID, err.Error()))
					return
				}

				msg += fmt.Sprintf("Валюта: %s Сумма: %.2f [%.2f]\n", key, value, value*coinPrice)
				usdSum += value * coinPrice
			}

			msg += fmt.Sprintf("Сумма USD: %.2f\n", usdSum)

			bot.Send(tgbotapi.NewMessage(update.Message.Chat.ID, msg))
		default:
			bot.Send(tgbotapi.NewMessage(update.Message.Chat.ID, "Команда неизвестна"))
		}

		// msg := tgbotapi.NewMessage(update.Message.Chat.ID, update.Message.Text)
		// msg.ReplyToMessageID = update.Message.MessageID

		// bot.Send(msg)
	}

}
