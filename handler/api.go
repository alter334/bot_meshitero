package handler

import (
	"context"
	"fmt"
	"os"

	"github.com/traPtitech/go-traq"
	traqwsbot "github.com/traPtitech/traq-ws-bot"
)

func SimplePost(bot *traqwsbot.Bot, channelID string, content string) (x string) {
	q, r, err := bot.API().
		MessageApi.
		PostMessage(context.Background(), channelID).
		PostMessageRequest(traq.PostMessageRequest{
			Content: content,
		}).
		Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
	return q.Id
}

func SimpleEdit(bot *traqwsbot.Bot, m string, s string) {
	bot.API().
		MessageApi.EditMessage(context.Background(), m).PostMessageRequest(traq.PostMessageRequest{
		Content: s,
	}).Execute()
}

func GetUserHome(bot *traqwsbot.Bot, postUserID string) string{
	userdetail, httpres, err := bot.API().UserApi.GetUser(context.Background(), postUserID).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", httpres)
	}
	return userdetail.GetHomeChannel()
}
