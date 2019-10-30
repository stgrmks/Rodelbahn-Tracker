package bot

import (
	"fmt"
	"github.com/nlopes/slack"
	"github.com/stgrmks/Rodelbahn-Tracker/internal/config"
	"github.com/stgrmks/Rodelbahn-Tracker/internal/crawler"
	"strings"
)

const (
	Version  = "VERSION"
	CrawlNow = "CRAWLNOW"
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
			commandHandler(userId, chanId, cmdSplit[1], rtm)

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

func commandHandler(userId string, chanId string, msg string, rtm *slack.RTM) {
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

	case CrawlNow:
		dbSess := crawler.DbSession{}
		dbSess.Connect(&MyConfig)
		crwl := crawler.Control{}
		crwl.Links = MyConfig.RbList
		rtm.SendMessage(rtm.NewOutgoingMessage(fmt.Sprintf("<@%s> Started Crawler...", userId), chanId))
		crwl.Start(&MyConfig)
		dbSess.Commit(crwl.Result)
		rtm.SendMessage(rtm.NewOutgoingMessage(fmt.Sprintf("<@%s> Crawler finished :)", userId), chanId))

	}

	log.Println("command: ", msg, " cmdSplit: ", cmdSplit)
}
