package internal

import (
	"fmt"
	"os"

	"github.com/slack-go/slack"
)

var (
	slackToken  = os.Getenv("SLACK_TOKEN")
	channelName = os.Getenv("CHANNEL_NAME")
)

func SendMessage(text string) {
	fmt.Println("Sending message to channel: " + text)
	api := slack.New(slackToken)
	_, _, err := api.PostMessage(channelName, slack.MsgOptionText(text, false))
	if err != nil {
		fmt.Println("Error occurred during sending slack message:", err)
	}
}

func SendMessageToChannel(channel string, text string) {
	fmt.Println("Sending message to channel: " + text)
	api := slack.New(slackToken)
	_, _, err := api.PostMessage(channel, slack.MsgOptionText(text, false))
	if err != nil {
		fmt.Println("Error occurred during sending slack message:", err)
	}
}
