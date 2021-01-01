package slackbot

import (
	"fmt"
	"github.com/slack-go/slack"
	"github.com/stgrmks/Rodelbahn-Tracker/internal/config"
	"strings"
)

type Bot struct {
	Version           string
	Build             string
	Shutdown          chan bool
	StopPeriodicCrawl chan bool
	MyConfig          config.Config
	CommandMap        map[string]*Command
	api               *slack.Client
	rtm               *slack.RTM
}

func (b *Bot) StartBot() {
	b.api = slack.New(b.MyConfig.SlackBotToken, slack.OptionDebug(true))
	b.rtm = b.api.NewRTM()
	go b.rtm.ManageConnection()

	for msg := range b.rtm.IncomingEvents {
		log.Debug("Event Received: ")
		switch ev := msg.Data.(type) {

		case *slack.ConnectedEvent:
			log.Debugln("Infos:", ev.Info)
			log.Debugln("Connection counter:", ev.ConnectionCount)
			b.rtm.SendMessage(b.rtm.NewOutgoingMessage("Hello world", "CFVJPFU2K"))

		case *slack.MessageEvent:
			bot := b.rtm.GetInfo().User.ID
			user := ev.User
			msg := ev.Text
			botTagInMsg := fmt.Sprintf("<@%s>", bot)

			if (bot == user) || (!strings.Contains(msg, botTagInMsg)) {
				log.Debugln("Msg from slackbot or slackbot was not addressed directly")
				continue
			}
			log.Debugln("bot: ", bot, "user: ", user, " Message: ", msg, " BotIdentifierInMsg: ", botTagInMsg)
			channel := ev.Channel
			msgSlice, err := msgSplitAndValidate(BiggerEqual, 2, msg, botTagInMsg)
			if err != nil {
				returnMsg := fmt.Sprintf("Sorry <@%s>, %s :(", user, err.Error())
				b.rtm.SendMessage(b.rtm.NewOutgoingMessage(returnMsg, channel))
			}
			b.invokeCommand(user, channel, msgSlice[1])

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

func (b *Bot) invokeCommand(user string, channel string, msg string) {
	var userParams []*Param
	msgSlice, err := msgSplitAndValidate(BiggerEqual, 2, msg, " ")
	if err != nil {
		log.Errorln("Error while splitting msg: ", err)
		b.rtm.SendMessage(b.rtm.NewOutgoingMessage(fmt.Sprintf("<@%s> Could not understand the command :(", user), channel))
		return
	}
	cmdString := msgSlice[1]
	cmd, ok := b.CommandMap[cmdString]
	if !ok {
		log.Errorf("Command %s does not exist.", cmdString)
		b.rtm.SendMessage(b.rtm.NewOutgoingMessage(fmt.Sprintf("<@%s> Command %s does not exist!", user, cmdString), channel))
		return
	}
	if len(msgSlice) > 2 {
		userParams = cmd.validateParams(msgSlice[2:])
	}
	cmd.execute(userParams, user, channel, b)
}
