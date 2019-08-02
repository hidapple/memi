package main

import (
	"os"

	"github.com/nlopes/slack"
)

func main() {
	c := slack.New(os.Getenv("SLACK_TOKEN"))
	bot := &SlackBot{client: c}
	bot.OpenRTM()
}
