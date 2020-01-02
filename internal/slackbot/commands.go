package slackbot

import (
	"fmt"
	"github.com/nlopes/slack"
	"github.com/robfig/cron"
	"github.com/stgrmks/Rodelbahn-Tracker/internal/crawler"
	"github.com/stgrmks/Rodelbahn-Tracker/internal/logger"
	"gopkg.in/mgo.v2/bson"
	"reflect"
	"strings"
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
	ParamMap    map[string]*Param
	Active      bool
}

func (c *Command) validateParams(msg []string) []*Param {
	var paramList []*Param
	var paramStruct *Param
	var ok bool

	for _, paramString := range msg {

		if strings.Contains(paramString, "::") {
			// we assume format param::value
			paramSlice, err := msgSplitAndValidate(Equal, 2, paramString, "::")
			if err != nil {
				log.Errorln("Error while splitting msg: ", err)
				continue
			}
			paramStruct, ok = c.ParamMap[paramSlice[0]]
			if !ok {
				log.Debug("Param %s does not exist.", paramString)
				continue
			}
			paramStruct.Value = paramSlice[1]
			paramList = append(paramList, paramStruct)

		} else {
			//	we assume format value which should be flag only
			paramStruct, ok := c.ParamMap[paramString]
			if !ok {
				log.Debug("Param %s does not exist.", paramString)
				continue
			}
			paramList = append(paramList, paramStruct)
		}
	}
	return paramList
}

func (c *Command) execute(ps []*Param, user string, channel string, b *Bot) {
	switch c.Name {
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
		c.Active = true
		startPeriodicCrawler(user, channel, ps, b)
		break

	case "stopPeriodicCrawl":
		if !b.CommandMap["periodicCrawl"].Active {
			b.rtm.SendMessage(b.rtm.NewOutgoingMessage(fmt.Sprintf("<@%s> Periodic crawl is not activated!", user), channel))
			break
		}
		b.CommandMap["periodicCrawl"].Active = false
		b.StopPeriodicCrawl <- true
		break

	case "shutdown":
		b.Shutdown <- true
		b.rtm.SendMessage(b.rtm.NewOutgoingMessage(fmt.Sprintf("<@%s> Going to sleep. Bye!", user), channel))
		break

	case "lastEntries":
		getLastEntries(user, channel, ps, b)
		break

	case "help":
		sendHelpMsg(user, channel, b)
		break
	}

}

func sendHelpMsg(user string, channel string, b *Bot) {
	msg := fmt.Sprintf("<@%s> Available Commands: \n", user)
	for _, cmd := range b.CommandMap {
		log.Infoln(cmd.Name)
		msg += fmt.Sprintf("%s: %s", cmd.Name, cmd.Description)
		if cmd.ParamMap != nil {
			msg += "Params: "
			for _, param := range cmd.ParamMap {
				msg += fmt.Sprintf("%s", param.Name)
			}
		}
		msg += "\n"
	}
	b.rtm.SendMessage(b.rtm.NewOutgoingMessage(msg, channel))
}

func startPeriodicCrawler(user, channel string, ps []*Param, b *Bot) {

	// cron setup
	cr := cron.New()
	err := cr.AddFunc(b.MyConfig.Cron, func() {
		startCrawler(user, channel, []*Param{}, b)
	})
	if err != nil {
		log.Errorln("Failed to add Crawl function to cron service", err)
		return
	}
	cr.Start()
	log.Debug("Started cron service.")
	b.rtm.SendMessage(b.rtm.NewOutgoingMessage(fmt.Sprintf("<@%s> Starting periodic crawl with cron-pattern: %s", user, b.MyConfig.Cron), channel))

	go stopPeriodicCrawlListener(user, channel, cr, b)
}

func stopPeriodicCrawlListener(user, channel string, cr *cron.Cron, b *Bot) {
	// TODO: move this as parameter into periodicCrawl
	// waiting for kill signal
	for {
		<-b.StopPeriodicCrawl
		cr.Stop()
		logger.Logger.Info("Stopped periodical crawl.")
		b.rtm.SendMessage(b.rtm.NewOutgoingMessage(fmt.Sprintf("<@%s> Stopped periodic crawl.", user), channel))
	}
}

func startCrawler(user, channel string, ps []*Param, b *Bot) {
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

func getLastEntries(user string, channel string, ps []*Param, b *Bot) {
	var result []crawler.RbData
	var msg string
	dbSess := crawler.DbSession{}
	dbSess.Connect(&b.MyConfig)

	var query bson.M = nil
	for _, p := range ps {
		if p.Name == "location" {
			query = bson.M{"location": p.Value}
		}
	}

	err := dbSess.Collection.Find(query).Limit(3).All(&result)
	if err != nil {
		log.Fatal("Could not retrieve the last entries!")
		msg = fmt.Sprintf("<@%s> Could not retrieve the last entries :(", user)
		b.rtm.SendMessage(b.rtm.NewOutgoingMessage(msg, channel))
	}

	if len(result) > 0 {
		var attachment slack.Attachment
		for _, r := range result {
			attachment = createAttachement(r, attachment)
			channelID, timestamp, err := b.api.PostMessage(channel, slack.MsgOptionText(fmt.Sprintf("<@%s>", user), false), slack.MsgOptionAttachments(attachment))
			if err != nil {
				log.Errorf("%s\n", err)
			}
			log.Debugf("Message successfully sent to channel %s at %s", channelID, timestamp)
		}
	} else {
		msg = fmt.Sprintf("<@%s> No entries :(", user)
		b.rtm.SendMessage(b.rtm.NewOutgoingMessage(msg, channel))
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
