package slackbot

import (
	"fmt"
	"github.com/nlopes/slack"
	"github.com/robfig/cron"
	"github.com/stgrmks/Rodelbahn-Tracker/internal/crawler"
	"github.com/stgrmks/Rodelbahn-Tracker/internal/logger"
	"reflect"
)

type Param struct {
	Name     string
	Value    string
	Optional bool
	Flag     bool
}

type Command struct {
	Name        string
	Description string
	ParamMap    map[string]Param
}

func (c *Command) validateParams(msg []string) []*Param {
	paramList := []*Param{}

	for _, param := range msg {
		paramStruct, ok := c.ParamMap[param]
		if !ok {
			log.Debugf("Param %s does not exist.", param)
			break
		}

		// just giving pointer is much faster. doesnt matter much for small slices though
		paramList = append(paramList, &paramStruct)
	}
	return paramList
}

func (c *Command) execute(command Command, ps []*Param, user string, channel string, b *Bot) {
	switch command.Name {
	case "version":
		sendVersionBuildMsg(user, channel, b)
		break

	case "showConfig":
		sendCurrentConfig(user, channel, b)
		break

	case "crawlNow":
		startCrawler(user, channel, ps, b)
		break

	case "periodicCrawl":
		startPeriodicCrawler(user, channel, ps, b)
		break

	}
}

func startPeriodicCrawler(user string, channel string, ps []*Param, b *Bot) {
	// cron setup
	c := cron.New()
	err := c.AddFunc(b.MyConfig.Cron, func() {
		startCrawler(user, channel, []*Param{}, b)
	})
	if err != nil {
		log.Errorln("Failed to add Crawl function to cron service", err)
		return
	}
	c.Start()
	log.Debug("Started cron service.")
	b.rtm.SendMessage(b.rtm.NewOutgoingMessage(fmt.Sprintf("<@%s> Starting Periodic Crawl with Cron: %s", user, b.MyConfig.Cron), channel))

	// waiting for kill signal
	<-b.StopPeriodicCrawl
	logger.Logger.Info("Stopping Periodical Crawl.")
}

func startCrawler(user string, channel string, ps []*Param, b *Bot) {
	silent := false
	for _, p := range ps {
		if p.Name == "silent" {
			silent = p.Flag
		}
	}
	var attachment slack.Attachment
	dbSess := crawler.DbSession{}
	dbSess.Connect(&b.MyConfig)
	crwl := crawler.Control{}
	crwl.Links = b.MyConfig.RbList
	crwl.Start(&b.MyConfig)
	newEntries := dbSess.Commit(crwl.Result)
	b.rtm.SendMessage(b.rtm.NewOutgoingMessage(fmt.Sprintf("<@%s> Crawler finished successfully in %s with %d new entries!", user, crwl.EndTime.Sub(crwl.StartTime), len(newEntries)), channel))
	if b.MyConfig.Notify && !silent {
		for _, entry := range newEntries {
			attachment = createAttachement(entry, attachment)
			channelID, timestamp, err := b.api.PostMessage(channel, slack.MsgOptionText(fmt.Sprintf("<@%s>", user), false), slack.MsgOptionAttachments(attachment))
			if err != nil {
				log.Errorf("%s\n", err)
			}
			log.Debugf("Message successfully sent to channel %s at %s", channelID, timestamp)
		}
	}
}

func sendCurrentConfig(user string, channel string, b *Bot) {
	key := reflect.TypeOf(b.MyConfig)
	value := reflect.ValueOf(b.MyConfig)
	msg := fmt.Sprintf("<@%s> Current Config\n", user)
	for i := 0; i < value.NumField(); i++ {
		msg = msg + fmt.Sprintf("%s - %s \n", key.Field(i).Name, value.Field(i))
	}
	b.rtm.SendMessage(b.rtm.NewOutgoingMessage(msg, channel))
}

func sendVersionBuildMsg(user string, channel string, b *Bot) {
	log.Print("Sending Version and Build Info")
	msg := fmt.Sprintf("<@%s> Version: %s Build: %s", user, b.Version, b.Build)
	b.rtm.SendMessage(b.rtm.NewOutgoingMessage(msg, channel))
}
