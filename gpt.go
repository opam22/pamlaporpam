package main

import (
	"context"
	"errors"
	"fmt"
	"io"
	"log"

	gogpt "github.com/sashabaranov/go-gpt3"
)

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
