package main

import (
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/joho/godotenv"
)

var (
	telegramtoken string
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

	openaitoken := os.Getenv("OPENAI_TOKEN")
	if openaitoken == "" {
		panic("missing telegram token")
	}

	temperature := os.Getenv("TEMPERATURE")
	if openaitoken == "" {
		panic("missing temperature")
	}
	ftemp, err := strconv.ParseFloat(temperature, 32)
	if err != nil {
		panic(err)
	}
	f32temp := float32(ftemp)

	maxToken := os.Getenv("MAX_TOKEN")
	if openaitoken == "" {
		panic("missing max_token")
	}
	imaxToken, err := strconv.Atoi(maxToken)
	if err != nil {
		fmt.Println(err)
		return
	}

	f, err := os.OpenFile("text.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Println(err)
	}
	defer f.Close()

	logger := log.New(f, "pamlaporpam", log.LstdFlags)
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
	gpt := newGPT(openaitoken, f32temp, imaxToken)

	// wait for messages
	updates := bot.GetUpdatesChan(u)
	for update := range updates {
		if update.Message != nil {
			if strings.Contains(update.Message.Text, "/start") {
				replyMsg := "Tanya aja nanti gw jawab"
				if err := reply(logger, bot, update.Message, replyMsg); err != nil {
					logger.Printf("error: %v\n", err.Error())
				}
				continue
			}

			g := update.Message.Chat.IsGroup()
			sg := update.Message.Chat.IsSuperGroup()
			if g || sg {
				if strings.Contains(update.Message.Text, "@PamLaporPamBot") {
					if strings.Contains(update.Message.Text, "abangku") || strings.Contains(update.Message.Text, "abangqu") {
						replyMsg := "iya ol..."
						if err := reply(logger, bot, update.Message, replyMsg); err != nil {
							logger.Printf("error: %v\n", err.Error())
						}
					} else {
						// this prompt require us to call openai gpt

						// send feedback first to user
						if err := reply(logger, bot, update.Message, "oke sebentar..."); err != nil {
							logger.Printf("error: %v\n", err.Error())
						}

						replyMsg, err := gpt.chatCompletion(update.Message.Text)
						if err != nil || replyMsg == "" {
							// gpt3 fail, send message to user
							if err := reply(logger, bot, update.Message, "ntar aja dah lagi puyeng gua..."); err != nil {
								logger.Printf("error: %v\n", err.Error())
							}
							logger.Printf("error: %v\n", err.Error())
						}
						if err := reply(logger, bot, update.Message, replyMsg); err != nil {
							logger.Printf("error: %v\n", err.Error())
						}
					}

				}
			} else {
				// this prompt require us to call openai gpt

				// send feedback first to user
				if err := reply(logger, bot, update.Message, "oke sebentar..."); err != nil {
					logger.Printf("error: %v\n", err.Error())
				}

				replyMsg, err := gpt.chatCompletion(update.Message.Text)
				if err != nil || replyMsg == "" {
					// gpt3 fail, send message to user
					if err := reply(logger, bot, update.Message, "ntar aja dah lagi puyeng gua..."); err != nil {
						logger.Printf("error: %v\n", err.Error())
					}
					logger.Printf("error: %v\n", err.Error())
				}

				if err := reply(logger, bot, update.Message, replyMsg); err != nil {
					logger.Printf("error: %v\n", err.Error())
				}
			}
		}
	}
}

func reply(logger *log.Logger, bot *tgbotapi.BotAPI, message *tgbotapi.Message, replyMsg string) error {
	r := tgbotapi.NewMessage(message.Chat.ID, replyMsg)
	r.ReplyToMessageID = message.MessageID
	if _, err := bot.Send(r); err != nil {
		return err
	}

	// logging
	logger.Println("===================================")
	logger.Printf("From: %s", message.From.UserName)
	logger.Printf("Message: %s", message.Text)
	logger.Printf("Reply: %s END", replyMsg)
	logger.Println("===================================")

	return nil
}
