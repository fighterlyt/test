package main

import (
	"log"
	"gopkg.in/telegram-bot-api.v4"
)

func main() {
	bot, err := tgbotapi.NewBotAPI("485828201:AAE1uJb-eB0NEyXgcWQGb3L953_RVudBbQ0")
	if err != nil {
		log.Panic(err)
	}

	bot.Debug = true

	log.Printf("Authorized on account %s", bot.Self.UserName)



	bot.Send(tgbotapi.MessageConfig{
		ParseMode:tgbotapi.ModeMarkdown,
		BaseChat:tgbotapi.BaseChat{
			ChatID:-1001199430805,
		},
		Text:"*test*",

	})

	bot.Send(tgbotapi.NewMessage(-1001199430805,`<b>test</b>`))
	//u := tgbotapi.NewUpdate(0)
	//u.Timeout = 60
	//
	//
	//
	//updates, err := bot.GetUpdatesChan(u)
	//
	//
	//
	//for update := range updates {
	//	if update.Message == nil {
	//		continue
	//	}
	//
	//	log.Printf("[%s] %s", update.Message.From.UserName, update.Message.Text)
	//
	//	msg := tgbotapi.NewMessage(update.Message.Chat.ID, update.Message.Text)
	//	msg.ReplyToMessageID = update.Message.MessageID
	//
	//	bot.Send(msg)
	//}
}