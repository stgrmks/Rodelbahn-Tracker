package main

import (
	"github.com/robfig/cron"
	"github.com/sbstjn/hanu"
)

var (
	killPeriodicCrawl = make(chan bool)
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
		killMain <- true
	})
	slack.Command("crawlNow", HandleCrawlNow)
	slack.Command("periodicCrawl", HandlePeriodicCrawl)
	slack.Command("stopPeriodicCrawl", func(conv hanu.ConversationInterface) {
		conv.Reply("Stopping Periodic Crawl")
		killPeriodicCrawl <- true
	})
	slack.Command("changeCron <Cron-Pattern>", HandleChangeCron)

	slack.Listen()
}

func HandleCrawlNow(conv hanu.ConversationInterface) {
	conv.Reply("Executing Crawler now...")

	ActiveCrawler := RunCrawler()
	conv.Reply("Crawler finished after %s", ActiveCrawler.EndTime.Sub(ActiveCrawler.StartTime).String())

}

func RunCrawler() CrawlerControl {
	// Establish DB connection
	ActiveDbSession := DbSession{}
	ActiveDbSession.Connect(&MyConfig)
	// Start Crawler
	ActiveCrawler := CrawlerControl{}
	ActiveCrawler.Links = MyConfig.RbList
	ActiveCrawler.Start(&MyConfig)
	// Insert crawler result to DB
	ActiveDbSession.Commit(ActiveCrawler.Result)
	return ActiveCrawler
}

func HandlePeriodicCrawl(conv hanu.ConversationInterface) {
	conv.Reply("Started periodic crawler with cron: %s", MyConfig.Cron)

	// cron setup
	c := cron.New()
	c.AddFunc(MyConfig.Cron, func() {
		_ = RunCrawler()
	})
	log.Infof("Periodical Crawl initiated: %s", MyConfig.Cron)
	c.Start()

	// waiting for kill signal
	<-killPeriodicCrawl
	log.Info("Stopping Periodical Crawl.")
}

func HandleChangeCron(conv hanu.ConversationInterface) {
	conv.Reply("TO BE IMPLEMENTED!")
}
