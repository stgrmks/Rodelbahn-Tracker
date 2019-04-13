package cmd

import (
	"fmt"
	"github.com/sbstjn/hanu"
	"github.com/spf13/cobra"
	"gopkg.in/mgo.v2/bson"
	"log"
	"os"
	"os/signal"
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

	slack.Command("lastEntry <Rodelbahn-Name>", HandleLastEntry)
	slack.Command("allEntries <Rodelbahn-Name>", HandleAllEntries)
	slack.Command("changeCron <Cron-Pattern>", HandleChangeCron)
	slack.Command("crawlNow", HandleCrawlNow)
	slack.Command("initCrawl", HandleInitCrawl)
	slack.Command("periodicCrawl", HandlePeriodicCrawl)
	slack.Command("version", func(conv hanu.ConversationInterface) {
		conv.Reply("Version: %s Build: %s", Version, Build)
	})
	slack.Command("ls", HandleLS)

	slack.Listen()

	sig := make(chan os.Signal)
	signal.Notify(sig, os.Interrupt, os.Kill)
	<-sig
}

func HandleLastEntry(conv hanu.ConversationInterface) {
	//query := ""
	results := RbData{}
	location, err := conv.String("Rodelbahn-Name")
	if err != nil {
		log.Fatal(err)
	}
	config.ActiveCollection.Find(bson.M{"location": location}).Sort("time").One(&results)
	fmt.Println(results)
	conv.Reply("TO BE IMPLEMENTED!")
}

func HandleAllEntries(conv hanu.ConversationInterface) {
	query := ""
	results := []RbData{}
	err := config.ActiveCollection.Find(query).Sort("time").All(&results)
	if err != nil {
		log.Println("Query failed: %s", err)
	}
	conv.Reply("TO BE IMPLEMENTED!")
}

func HandleChangeCron(conv hanu.ConversationInterface) {
	conv.Reply("TO BE IMPLEMENTED!")
}

func HandleCrawlNow(conv hanu.ConversationInterface) {
	conv.Reply("Executing Crawler now...")
	RunStartCrawler()
}

func HandleInitCrawl(conv hanu.ConversationInterface) {
	RunInitCrawl()
	conv.Reply("Updated database without notifications...")
	RunStartCrawler()
}

func HandlePeriodicCrawl(conv hanu.ConversationInterface) {
	go RunPeriodicCrawler()
	conv.Reply("Started periodic crawler...")
}

func HandleLS(conv hanu.ConversationInterface) {
	RunInitDB()
	result := []string{}
	if err := config.ActiveCollection.Find(nil).Distinct("location", &result); err != nil {
		log.Println(err)
	}
	fmt.Println(result)
}
