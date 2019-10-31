package bot

import (
	"fmt"
	"github.com/nlopes/slack"
	"github.com/robfig/cron"
	"github.com/stgrmks/Rodelbahn-Tracker/internal/config"
	"github.com/stgrmks/Rodelbahn-Tracker/internal/crawler"
	"github.com/stgrmks/Rodelbahn-Tracker/internal/logger"
	"reflect"
	"strings"
)

const (
	Version       = "VERSION"
	CrawlNow      = "CRAWLNOW"
	ShowConfig    = "SHOWCONFIG"
	PeriodicCrawl = "PERIODICCRAWL"
)

var (
	VERSION           = "0.2.0"
	BUILD             = "0.2.0"
	KillBot           = make(chan bool)
	KillPeriodicCrawl = make(chan bool)
	MyConfig          config.Config
)

func StartBot() {
	api := slack.New(MyConfig.SlackBotToken, slack.OptionDebug(true))
	rtm := api.NewRTM()
	go rtm.ManageConnection()

	for msg := range rtm.IncomingEvents {
		log.Debug("Event Received: ")
		switch ev := msg.Data.(type) {

		case *slack.ConnectedEvent:
			log.Debugln("Infos:", ev.Info)
			log.Debugln("Connection counter:", ev.ConnectionCount)
			rtm.SendMessage(rtm.NewOutgoingMessage("Hello world", "C2147483705"))

		case *slack.MessageEvent:
			botId := rtm.GetInfo().User.ID
			userId := ev.User
			msg := ev.Text
			botTagInMsg := fmt.Sprintf("<@%s>", botId)

			if (botId == userId) || (!strings.Contains(msg, botTagInMsg)) {
				log.Debugln("Msg from bot or bot was not addressed directly")
				continue
			}
			log.Debugln("botId: ", botId, "userId: ", userId, " Message: ", msg, " BotIdentifierInMsg: ", botTagInMsg)
			chanId := ev.Channel
			cmdSplit, err := msgSplit(Equal, 2, msg, botTagInMsg)
			if err != nil {
				returnMsg := fmt.Sprintf("Sorry <@%s>, %s :(", userId, err.Error())
				rtm.SendMessage(rtm.NewOutgoingMessage(returnMsg, chanId))
			}
			commandHandler(userId, chanId, cmdSplit[1], rtm, api)

		case *slack.PresenceChangeEvent:
			log.Debugf("Presence Change: %v\n", ev)

		case *slack.LatencyReport:
			log.Debugf("Current latency: %v\n", ev.Value)

		case *slack.RTMError:
			log.Debugf("Error: %s\n", ev.Error())

		case *slack.InvalidAuthEvent:
			log.Debugf("Invalid credentials")
		}
	}

}

func commandHandler(userId string, chanId string, msg string, rtm *slack.RTM, api *slack.Client) {
	returnMsg := ""
	cmdSplit, err := msgSplit(Bigger, 0, msg, " ")
	if err != nil {
		returnMsg = fmt.Sprintf("Sorry <@%s>, %s :(", userId, err.Error())
		rtm.SendMessage(rtm.NewOutgoingMessage(returnMsg, chanId))
	}
	cmd := strings.ToUpper(cmdSplit[1])
	fmt.Println(cmd)
	switch cmd {
	case Version:
		log.Print("Sending Version and Build Info")
		returnMsg = fmt.Sprintf("<@%s> Version: %s Build: %s", userId, VERSION, BUILD)
		rtm.SendMessage(rtm.NewOutgoingMessage(returnMsg, chanId))
		break

	case ShowConfig:
		key := reflect.TypeOf(MyConfig)
		value := reflect.ValueOf(MyConfig)
		rtnStr := fmt.Sprintf("<@%s> Current Config\n", userId)
		for i := 0; i < value.NumField(); i++ {
			rtnStr = rtnStr + fmt.Sprintf("%s - %s \n", key.Field(i).Name, value.Field(i))
		}
		rtm.SendMessage(rtm.NewOutgoingMessage(rtnStr, chanId))
		break

	case CrawlNow:
		startCrawler(rtm, userId, chanId, api)
		break

	case PeriodicCrawl:
		// cron setup
		c := cron.New()
		err := c.AddFunc(MyConfig.Cron, func() {
			startCrawler(rtm, userId, chanId, api)
		})
		if err != nil {
			log.Errorln("Failed to add Crawl function to cron service", err)
			return
		}
		c.Start()
		log.Debug("Started cron service.")
		rtm.SendMessage(rtm.NewOutgoingMessage(fmt.Sprintf("<@%s> Starting Periodic Crawl with Cron: %s", userId, MyConfig.Cron), chanId))

		// waiting for kill signal
		<-KillPeriodicCrawl
		logger.Logger.Info("Stopping Periodical Crawl.")
		break

	}

	log.Println("command: ", msg, " cmdSplit: ", cmdSplit)
}

func startCrawler(rtm *slack.RTM, userId string, chanId string, api *slack.Client) {
	var attachment slack.Attachment
	dbSess := crawler.DbSession{}
	dbSess.Connect(&MyConfig)
	crwl := crawler.Control{}
	crwl.Links = MyConfig.RbList
	crwl.Start(&MyConfig)
	newEntries := dbSess.Commit(crwl.Result)
	rtm.SendMessage(rtm.NewOutgoingMessage(fmt.Sprintf("<@%s> Crawler finished successfully in %s with %d new entries!", userId, crwl.EndTime.Sub(crwl.StartTime), len(newEntries)), chanId))
	if MyConfig.Notify {
		for _, entry := range newEntries {
			preText := fmt.Sprintf("New Rating for %s!", entry.Location)
			text := fmt.Sprintf("Timestamp: %s\nUser: %s\nRating: %s\nComment: %s", entry.Time.Format("2006-01-02"), entry.User, entry.Rating, entry.Comment)
			attachment = slack.Attachment{
				Title:      entry.Location,
				Pretext:    preText,
				Text:       text,
				TitleLink:  entry.Link,
				MarkdownIn: []string{"text", "pretext"},
			}
			channelID, timestamp, err := api.PostMessage(chanId, slack.MsgOptionText(fmt.Sprintf("<@%s>", userId), false), slack.MsgOptionAttachments(attachment))
			if err != nil {
				log.Errorf("%s\n", err)
			}
			log.Debugf("Message successfully sent to channel %s at %s", channelID, timestamp)
		}
	}
}
