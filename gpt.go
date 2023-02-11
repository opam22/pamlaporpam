package main

import (
	"context"
	"errors"
	"fmt"
	"io"

	gogpt "github.com/sashabaranov/go-gpt3"
)

type GPT struct {
	c   *gogpt.Client
	ctx context.Context
}

func newGPT() *GPT {
	c := gogpt.NewClient(openaitoken)
	return &GPT{
		c:   c,
		ctx: context.Background(),
	}
}

func (g *GPT) call(msg string) (string, error) {
	req := gogpt.CompletionRequest{
		Model:       gogpt.GPT3TextDavinci003,
		MaxTokens:   500,
		Prompt:      msg,
		Stream:      true,
		Temperature: 0,
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

		for _, ch := range response.Choices {
			reply = fmt.Sprintf("%s%s", reply, ch.Text)
		}
	}
}
