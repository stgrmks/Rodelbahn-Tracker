package rbtracker

import (
	"github.com/robfig/cron"
	"github.com/sbstjn/hanu"
	"github.com/stgrmks/Rodelbahn-Tracker/internal/config"
	"github.com/stgrmks/Rodelbahn-Tracker/internal/logger"
	"reflect"
)

var (
	VERSION           = "0.2.0"
	BUILD             = "0.2.0"
	KillBot           = make(chan bool)
	KillPeriodicCrawl = make(chan bool)
	MyConfig          config.Config
)

func StartBot() {

	slack, err := hanu.New(MyConfig.SlackBotToken)
	if err != nil {
		logger.Logger.Fatal(err)
	}

	slack.Command("version", func(conv hanu.ConversationInterface) {
		conv.Reply("Version: %s Build: %s", VERSION, BUILD)
	})
	slack.Command("kill", func(conv hanu.ConversationInterface) {
		conv.Reply("bye bye")
		KillBot <- true
	})
	slack.Command("showConfig", HandleShowConfig)
	slack.Command("crawlNow", HandleCrawlNow)
	slack.Command("startPeriodicCrawl", HandlePeriodicCrawl)
	slack.Command("stopPeriodicCrawl", func(conv hanu.ConversationInterface) {
		conv.Reply("Stopping Periodic Crawl")
		KillPeriodicCrawl <- true
	})
	slack.Command("changeCron <Cron-Pattern>", HandleChangeCron)

	slack.Listen()
}

func HandleShowConfig(conv hanu.ConversationInterface) {
	key := reflect.TypeOf(MyConfig)
	value := reflect.ValueOf(MyConfig)

	for i := 0; i < value.NumField(); i++ {
		conv.Reply("%s : %s", key.Field(i).Name, value.Field(i))
	}
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
	logger.Logger.Infof("Periodical Crawl initiated: %s", MyConfig.Cron)
	c.Start()

	// waiting for kill signal
	<-KillPeriodicCrawl
	logger.Logger.Info("Stopping Periodical Crawl.")
}

func HandleChangeCron(conv hanu.ConversationInterface) {
	newCron, _ := conv.String("Cron-Pattern")
	MyConfig.Cron = newCron
	logger.Logger.Infof("Changing Cron-Pattern to: %s", MyConfig.Cron)
	conv.Reply("Changing Cron-Pattern to: `%s`", MyConfig.Cron)
}
