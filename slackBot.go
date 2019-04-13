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
		MainIsDone <- true
	})

	slack.Command("testDB", HandleCrawlNow)

	slack.Listen()
}

func HandleCrawlNow(conv hanu.ConversationInterface) {
	conv.Reply("Executing Crawler now...")

	// Establish DB connection
	ActiveDbSession := DbSession{}
	ActiveDbSession.Connect(&MyConfig)

	// Start Crawler
	ActiveCrawler := CrawlerControl{}
	ActiveCrawler.Links = MyConfig.RbList
	ActiveCrawler.Start(&MyConfig)

	// Insert crawler result to DB

}
