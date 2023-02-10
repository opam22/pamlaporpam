package main

import (
	"log"
	"os"
	"strings"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/joho/godotenv"
)

var (
	telegramtoken, openaitoken string
)

func main() {
	err := godotenv.Load("config.env")
	if err != nil {
		log.Fatalf("Some error occured. Err: %s", err)
	}

	telegramtoken = os.Getenv("TELEGERAM_TOKEN")
	if telegramtoken == "" {
		panic("missing telegram token")
	}

	openaitoken = os.Getenv("OPENAI_TOKEN")
	if openaitoken == "" {
		panic("missing telegram token")
	}

	f, err := os.OpenFile("text.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Println(err)
	}
	defer f.Close()

	logger := log.New(f, "prefix", log.LstdFlags)
	logger.Println("===================================")

	bot, err := tgbotapi.NewBotAPI(telegramtoken)
	if err != nil {
		log.Panic(err)
	}

	bot.Debug = true
	log.Printf("Authorized on account %s", bot.Self.UserName)
	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	// create instance for gpt
	gpt := newGPT()

	// wait for messages
	updates := bot.GetUpdatesChan(u)
	for update := range updates {
		if update.Message != nil {
			logger.Println("===================================")
			logger.Printf("From: %s", update.Message.From.UserName)
			logger.Printf("Message: %s", update.Message.Text)

			if strings.Contains(update.Message.Text, "/start") {
				reply := "Tanya aja nanti gw jawab"
				r := tgbotapi.NewMessage(update.Message.Chat.ID, reply)
				r.ReplyToMessageID = update.Message.MessageID
				bot.Send(r)
				logger.Printf("Reply: %s END", reply)
				logger.Println("===================================")
				continue
			}

			g := update.Message.Chat.IsGroup()
			sg := update.Message.Chat.IsSuperGroup()
			if g || sg {
				if strings.Contains(update.Message.Text, "@PamLaporPamBot") {
					if strings.Contains(update.Message.Text, "abangku") || strings.Contains(update.Message.Text, "abangqu") {
						reply := "iya ol..."
						r := tgbotapi.NewMessage(update.Message.Chat.ID, reply)
						r.ReplyToMessageID = update.Message.MessageID
						bot.Send(r)
						logger.Printf("Reply: %s END", reply)
						logger.Println("===================================")
					} else {
						reply, err := gpt.call(update.Message.Text)
						if err != nil {
							log.Println("callgpt error:", err.Error())
						}
						r := tgbotapi.NewMessage(update.Message.Chat.ID, reply)
						r.ReplyToMessageID = update.Message.MessageID
						bot.Send(r)
						logger.Printf("Reply: %s END", reply)
						logger.Println("===================================")
					}

				}
			} else {
				reply, err := gpt.call(update.Message.Text)
				if err != nil {
					log.Println("callgpt error:", err.Error())
				}
				r := tgbotapi.NewMessage(update.Message.Chat.ID, reply)
				r.ReplyToMessageID = update.Message.MessageID
				bot.Send(r)
				logger.Printf("Reply: %s END", reply)
				logger.Println("===================================")
			}
		}
	}
}
