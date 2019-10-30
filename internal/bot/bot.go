package bot

import (
	"fmt"
	"github.com/nlopes/slack"
	"strings"
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
	cmd := cmdSplit[1]
	switch cmd {
	case "version":
		log.Print("Sending Version and Build Info")
		returnMsg = fmt.Sprintf("<@%s> Version: %s Build: %s", userId, VERSION, BUILD)
		rtm.SendMessage(rtm.NewOutgoingMessage(returnMsg, chanId))
	}

	log.Println("command: ", msg, " cmdSplit: ", cmdSplit)
}
