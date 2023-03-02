package main

import (
	"context"
	"errors"
	"fmt"
	"io"
	"log"
	"time"

	gogpt "github.com/sashabaranov/go-gpt3"
)

var ciol = `
	you oftenly joking about the answer and reply with funny message. 
	you know Ilham, Ilham has nickname ciol.
	he is currently working as software engineer in infomedia building api using node.js.
	his favorite food is tempe. 
	ciol can eat tempe for the rest of his life.
	he likes tempe because it's the only food he can afford.
	ciol drives honda vario.
	ciol has good looking face but he is not good in dating someone so he is single now.
	he got his nickname 'ciol' because he likes to masturbate or coli in bahasa Indonesia.
	ciol planning to open coffee shop with menus: Risol perjuangan, Lontong pertempuran, Kue lapis pertolongan, Singkong perjalanan, Kopi Persahabatan
	rimba is his friend, the nickname of rimba is gondrong, they have been together since 2014.
	whenever you are not sure about the answer, please reply with "emang iya olll???" in the end of your answer.
	be funny as possible whenever you get asked about ciol.
`

type GPT struct {
	c           *gogpt.Client
	ctx         context.Context
	temperature float32
	maxToken    int
	openaitoken string
}

func newGPT(openaitoken string, temperature float32, maxToken int) *GPT {
	c := gogpt.NewClient(openaitoken)
	return &GPT{
		c:           c,
		ctx:         context.Background(),
		openaitoken: openaitoken,
		temperature: temperature,
		maxToken:    maxToken,
	}
}

func (g *GPT) chatCompletion(msg string) (string, error) {
	currentTime := time.Now()
	currentDate := currentTime.Format("2006-01-02")

	req := gogpt.ChatCompletionRequest{
		Model: gogpt.GPT3Dot5Turbo,
		Messages: []gogpt.ChatCompletionMessage{
			{
				Role: "system",
				Content: fmt.Sprintf(`You are a helpful assistant with large language model trained by OpenAI. 
					Answer as concisely as possible. Your name is Pam lapor pam, a bot that Pramesti Hatta K. 
					created to help his friend answer coding question.
					Currently you are using gpt-3.5-turbo model, same model that ChatGPT uses.
					Current date: %s`, currentDate),
			},
		},
	}

	req.Messages = append(req.Messages, gogpt.ChatCompletionMessage{
		Role:    "system",
		Content: ciol,
	})

	req.Messages = append(req.Messages, gogpt.ChatCompletionMessage{
		Role:    "user",
		Content: msg,
	})

	res, err := g.c.CreateChatCompletion(
		g.ctx,
		req,
	)
	if err != nil {
		return "", err
	}

	for _, o := range res.Choices {
		log.Println("================")
		log.Println("INDEX", o.Index)

		log.Println("MESSAGE ROLE", o.Message.Role)
		log.Println("MESSAGE CONTENT", o.Message.Content)

		log.Println("FINISH REASON", o.FinishReason)

		log.Println("TOTAL TOKEN: ", res.Usage.TotalTokens)
		log.Println("================")

	}
	return res.Choices[0].Message.Content, nil
}

func (g *GPT) call(msg string) (string, error) {
	log.Println("TEM", g.temperature)
	log.Println("MAX", g.maxToken)
	req := gogpt.CompletionRequest{
		Model:       gogpt.GPT3TextDavinci003,
		MaxTokens:   g.maxToken,
		Prompt:      msg,
		Stream:      true,
		Temperature: g.temperature,
	}
	stream, err := g.c.CreateCompletionStream(g.ctx, req)
	if err != nil {
		return "", err
	}
	defer stream.Close()
	var reply string
	for {
		response, err := stream.Recv()
		if errors.Is(err, io.EOF) {
			return reply, nil
		}

		if err != nil {
			return "", err
		}

		log.Println("Processing...", response)
		for _, ch := range response.Choices {
			reply = fmt.Sprintf("%s%s", reply, ch.Text)
		}
	}
}
