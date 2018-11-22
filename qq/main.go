package main

import (
	"github.com/catsworld/qq-bot-api"
	"log"
	"net/http"
)

func main() {
	bot, err := qqbotapi.NewBotAPI("MyCoolqHttpToken", "http://localhost:5700", "CQHTTP_SECRET")
	if err != nil {
		log.Fatal(err)
	}

	bot.Debug = true

	u := qqbotapi.NewWebhook("/event")
	u.PreloadUserInfo = true

	updates := bot.ListenForWebSocket(u)

	api:=qqbotapi.NewWebhook("/api")
	api.PreloadUserInfo=true


	go http.ListenAndServe("0.0.0.0:12345", nil)

	for update := range updates {
		if update.Message == nil {
			continue
		}

		log.Printf("[%s] %s", update.Message.From.String(), update.Message.Text)

		msg := qqbotapi.NewMessage(update.Message.Chat.ID, update.Message.Chat.Type, update.Message.Text)
		bot.Send(msg)
	}

}
