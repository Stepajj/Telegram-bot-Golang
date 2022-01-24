package main

import (
	"fmt"
	"log"
	"strconv"
	"strings"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type bnResponse struct {
	Price float64 `json:"price, string"`
	Code  int64   `json:"code"`
}

type wallet map[string]float64

var db = map[int64]wallet{}

func main() {

	bot, err := tgbotapi.NewBotAPI("2134013692:AAH69EXzGhEb_y5gQWKerpf7GycapnAxr1s")
	if err != nil {
		log.Panic(err)
	}

	bot.Debug = true

	log.Printf("Authorized on account %s", bot.Self.UserName)

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates := bot.GetUpdatesChan(u)

	for update := range updates {
		if update.Message == nil { // ignore any non-Message Updates

			continue

		}

		//bot.Send(tgbotapi.NewMessage(update.Message.Chat.ID, "АУУУУУ"))

		command := strings.Split(update.Message.Text, " ")

		switch command[0] {
		case "ADD":
			if len(command) != 3 {
				bot.Send(tgbotapi.NewMessage(update.Message.Chat.ID, "Параменты команды заданы не правильно, пример правильной команды: ADD BTC 12"))
				continue
			}
			amount, err := strconv.ParseFloat(command[2], 64)
			if err != nil {
				bot.Send(tgbotapi.NewMessage(update.Message.Chat.ID, err.Error()))
			}

			if _, ok := db[update.Message.Chat.ID]; !ok {
				db[update.Message.Chat.ID] = wallet{}
			}

			db[update.Message.Chat.ID][command[1]] += amount

			balanceText := fmt.Sprintf("Пополнение произошло успешно, ваш баланс: %f", db[update.Message.Chat.ID][command[1]])
			bot.Send(tgbotapi.NewMessage(update.Message.Chat.ID, balanceText))

			//bot.Send(tgbotapi.NewMessage(update.Message.Chat.ID, "Валюта добавлена" ))
		case "SUB":

			if len(command) != 3 {
				bot.Send(tgbotapi.NewMessage(update.Message.Chat.ID, "Параменты команды заданы не правильно, пример правильной команды: SUB EUR 50"))
				continue
			}
			amount, err := strconv.ParseFloat(command[2], 64)
			if err != nil {
				bot.Send(tgbotapi.NewMessage(update.Message.Chat.ID, err.Error()))
			}

			if _, ok := db[update.Message.Chat.ID]; !ok {
				continue
			}

			db[update.Message.Chat.ID][command[1]] -= amount

			balanceText := fmt.Sprintf("Вычитание произошло успешно, остаток: %f", db[update.Message.Chat.ID][command[1]])
			bot.Send(tgbotapi.NewMessage(update.Message.Chat.ID, balanceText))

		case "DEL":
			if len(command) != 2 {
				bot.Send(tgbotapi.NewMessage(update.Message.Chat.ID, "Параменты команды заданы не правильно, пример правильной команды: DEL BTC"))
				continue
			}

			delete(db[update.Message.Chat.ID], command[1])

		case "SHOW":
			msg := ""
			for kay, value := range db[update.Message.Chat.ID] {
				msg += fmt.Sprintf("%s: %f\n", kay, value)
			}
			bot.Send(tgbotapi.NewMessage(update.Message.Chat.ID, msg))
		case "LOVE":
			bot.Send(tgbotapi.NewMessage(update.Message.Chat.ID, "Я очень сильно тебя люблю бип бип"))
		default:
			bot.Send(tgbotapi.NewMessage(update.Message.Chat.ID, "Команда не найдена, доступные команды:  "))
			bot.Send(tgbotapi.NewMessage(update.Message.Chat.ID, "ADD - добавление показателя "))
			bot.Send(tgbotapi.NewMessage(update.Message.Chat.ID, "SUB - вычитание показателя"))
			bot.Send(tgbotapi.NewMessage(update.Message.Chat.ID, "DEL - удаление показателя"))
			bot.Send(tgbotapi.NewMessage(update.Message.Chat.ID, "SHOW - список всех показателей"))
		}

		//log.Printf("[%s] %s", update.Message.From.UserName, update.Message.Text)

		//msg := tgbotapi.NewMessage(update.Message.Chat.ID, command[0])

		//bot.Send(msg)

	}

}
