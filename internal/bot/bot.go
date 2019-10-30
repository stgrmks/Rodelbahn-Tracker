package bot

import (
	"fmt"
	"github.com/nlopes/slack"
	"github.com/stgrmks/Rodelbahn-Tracker/internal/config"
	"github.com/stgrmks/Rodelbahn-Tracker/internal/logger"
	"strings"
)

var log = logger.Logger.WithField("package", "bot")

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
			botIdentifierInMsg := fmt.Sprintf("<@%s>", botId)

			if (botId == userId) || (!strings.Contains(msg, botIdentifierInMsg)) {
				log.Debugln("Msg from bot or bot was not addressed directly")
				continue
			}

			log.Debugln("botId: ", botId, "userId: ", userId, " Message: ", msg, " BotIdentifierInMsg: ", botIdentifierInMsg)

			rtm.SendMessage(rtm.NewOutgoingMessage("IT'S BURNTTTT!", ev.Channel))

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
