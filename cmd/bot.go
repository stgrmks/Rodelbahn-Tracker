package cmd

import (
	"github.com/sbstjn/hanu"
	"github.com/spf13/cobra"
	"log"
)

var bot = &cobra.Command{
	Use:   "bot",
	Short: "Starts Bot",
	Long:  `Starts an interactive C&C Slackbot.'`,
	Run: func(cmd *cobra.Command, args []string) {
		RunBot()
	},
}

func RunBot() {

	slack, err := hanu.New(config.SlackBotToken)
	if err != nil {
		log.Fatal(err)
	}

	slack.Command("LastEntry <Rodelbahn-Name>", HandleLastEntry)
	slack.Command("AllEntries <Rodelbahn-Name>", HandleAllEntries)
	slack.Command("ChangeCron <Cron-Pattern>", HandleChangeCron)
	slack.Command("CrawlNow", HandleCrawlNow)
	slack.Command("Version", func(conv hanu.ConversationInterface) {
		conv.Reply("Version: %s Build: %s", Version, Build)
	})

	slack.Listen()
}

func HandleLastEntry(conv hanu.ConversationInterface) {
	conv.Reply("TO BE IMPLEMENTED!")
}

func HandleAllEntries(conv hanu.ConversationInterface) {
	conv.Reply("TO BE IMPLEMENTED!")
}

func HandleChangeCron(conv hanu.ConversationInterface) {
	conv.Reply("TO BE IMPLEMENTED!")
}

func HandleCrawlNow(conv hanu.ConversationInterface) {
	conv.Reply("Executing Crawler now...")
	RunStartCrawler()
}
