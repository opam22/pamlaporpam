package main

import (
	"context"
	"errors"
	"fmt"
	"io"
	"log"
	"os"
	"strings"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	gogpt "github.com/sashabaranov/go-gpt3"
)

const (
	telegramtoken = ""
	openaitoken   = ""
)

func main() {
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

	updates := bot.GetUpdatesChan(u)
	for update := range updates {
		if update.Message != nil {
			logger.Printf("From: %s", update.Message.From.UserName)
			logger.Printf("Message: %s", update.Message.Text)
			g := update.Message.Chat.IsGroup()
			if g {
				if strings.Contains(update.Message.Text, "@PamLaporPamBot") {
					reply, err := callgpt(update.Message.Text)
					if err != nil {
						log.Println("callgpt error:", err.Error())
					}
					r := tgbotapi.NewMessage(update.Message.Chat.ID, reply)
					r.ReplyToMessageID = update.Message.MessageID
					bot.Send(r)
					logger.Printf("Reply: %s END", reply)
					logger.Println("===================================")
				}
			} else {
				msg := tgbotapi.NewMessage(update.Message.Chat.ID, update.Message.Text)
				msg.ReplyToMessageID = update.Message.MessageID

				logger.Printf("Reply: %s END", msg.Text)
				logger.Println("===================================")
				bot.Send(msg)
			}

			// reply, err := callgpt(msg.Text)
			// if err != nil {
			// 	log.Println("callgpt error:", err.Error())
			// }
			// r := tgbotapi.NewMessage(update.Message.Chat.ID, reply)
			// r.ReplyToMessageID = update.Message.MessageID
			// bot.Send(r)
		}
	}
}

func callgpt(msg string) (string, error) {
	c := gogpt.NewClient(openaitoken)
	ctx := context.Background()

	req := gogpt.CompletionRequest{
		Model:     gogpt.GPT3TextDavinci003,
		MaxTokens: 500,
		Prompt:    msg,
		Stream:    true,
	}
	stream, err := c.CreateCompletionStream(ctx, req)
	if err != nil {
		return "", err
	}
	defer stream.Close()
	var reply string
	for {
		response, err := stream.Recv()
		if errors.Is(err, io.EOF) {
			// log.Println("Stream finished")
			// log.Println("reply: ", reply)
			return reply, nil
		}

		if err != nil {
			return "", err
		}

		for _, ch := range response.Choices {
			// log.Println("loading ...", ch.Text)
			reply = fmt.Sprintf("%s%s", reply, ch.Text)
		}
	}
}
