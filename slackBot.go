package main

import (
	"github.com/sbstjn/hanu"
)

func StartBot(c *Config) {

	slack, err := hanu.New(c.SlackBotToken)
	if err != nil {
		log.Fatal(err)
	}

	slack.Command("version", func(conv hanu.ConversationInterface) {
		conv.Reply("Version: %s Build: %s", VERSION, BUILD)
	})
	slack.Command("kill", func(conv hanu.ConversationInterface) {
		conv.Reply("bye bye")
		mainFinished <- true
	})

	slack.Listen()
}
