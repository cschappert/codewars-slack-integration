package main

import (
	"fmt"
	"os"

	"github.com/slack-go/slack"
)

func main() {
	// read access token from env var and create new client
	api := slack.New(os.Getenv("SLACK_TOKEN"))

	// get channel ID from env var. bot must already be invited to room.
	channelId := os.Getenv("CHANNEL_ID")

	// send message to channel by ID with variable options.
	// TODO: is it possible to send to all joined / invited channels?
	_, _, _, err := api.SendMessage(
		channelId,
		slack.MsgOptionText("Mama Mia!", false),
		slack.MsgOptionUsername("The Pizza Man"),
		slack.MsgOptionIconEmoji(":pizza:"),
	)

	if err != nil {
		fmt.Println(err)
	}
}
