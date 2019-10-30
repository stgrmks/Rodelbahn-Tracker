package bot

import (
	"errors"
	"fmt"
	"github.com/nlopes/slack"
	"github.com/stgrmks/Rodelbahn-Tracker/internal/config"
	"github.com/stgrmks/Rodelbahn-Tracker/internal/logger"
	"strings"
)

var log = logger.Logger.WithField("package", "bot")

const (
	Equal        = "EQUAL"
	Inequal      = "INEQUAL"
	Bigger       = "BIGGER"
	BiggerEqual  = "BIGGEREQUAL"
	Smaller      = "SMALLER"
	SmallerEqual = "SMALLEREQUAL"
	All          = "ALL"
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

func msgSplit(cmdIsToSplits string, splits int, msg string, substr string) ([]string, error) {
	cmdRaw := strings.Split(msg, substr)
	cmpResult := false
	switch cmdIsToSplits {

	case Equal:
		cmpResult = len(cmdRaw) == splits
		break

	case Inequal:
		cmpResult = len(cmdRaw) != splits
		break

	case Bigger:
		cmpResult = len(cmdRaw) > splits
		break

	case BiggerEqual:
		cmpResult = len(cmdRaw) >= splits
		break

	case Smaller:
		cmpResult = len(cmdRaw) < splits
		break

	case SmallerEqual:
		cmpResult = len(cmdRaw) <= splits
		break

	case All:
		cmpResult = true
		break
	}

	if !cmpResult {
		err := errors.New("unknown command")
		log.Errorln("Comparison: ", cmpResult, "CommandRaw Len: ", len(cmdRaw), " Splits: ", splits, " CommandRaw: ", cmdRaw, " Error: ", err)
		return cmdRaw, err
	}
	return cmdRaw, nil
}

func commandHandler(userId string, chanId string, msg string, rtm *slack.RTM) {
	cmdSplit, err := msgSplit(Bigger, 0, msg, " ")
	if err != nil {
		returnMsg := fmt.Sprintf("Sorry <@%s>, %s :(", userId, err.Error())
		rtm.SendMessage(rtm.NewOutgoingMessage(returnMsg, chanId))
	}

	log.Println("command: ", msg, " cmdSplit: ", cmdSplit)
}
