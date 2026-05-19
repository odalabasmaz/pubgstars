package internal

import (
	"fmt"
	"github.com/bluele/slack"
	"os"
)

var (
	token       = os.Getenv("SLACK_TOKEN")
	channelName = os.Getenv("CHANNEL_NAME") //#customer-requests
)

func SendMessage(text string) {
	fmt.Println("Sending message to channel: " + text)

	api := slack.New(token)
	err := api.ChatPostMessage(channelName, text, nil)
	if err != nil {
		fmt.Println("Error occurred during sending slack message: ")
		fmt.Println(err)
	}
}

func SendMessageToChannel(channel string, text string) {
	fmt.Println("Sending message to channel: " + text)

	api := slack.New(token)
	err := api.ChatPostMessage(channel, text, nil)
	if err != nil {
		fmt.Println("Error occurred during sending slack message: ")
		fmt.Println(err)
	}
}
