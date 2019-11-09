package slackbot

import (
	"errors"
	"fmt"
	"github.com/nlopes/slack"
	"github.com/stgrmks/Rodelbahn-Tracker/internal/crawler"
	"github.com/stgrmks/Rodelbahn-Tracker/internal/logger"
	"strings"
)

var log = logger.Logger.WithField("package", "slackbot")

const (
	Equal        = "EQUAL"
	Inequal      = "INEQUAL"
	Bigger       = "BIGGER"
	BiggerEqual  = "BIGGEREQUAL"
	Smaller      = "SMALLER"
	SmallerEqual = "SMALLEREQUAL"
	All          = "ALL"
)

func msgSplitAndValidate(cmdIsToSplits string, splits int, msg string, substr string) ([]string, error) {
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

func createAttachement(entry crawler.RbData, attachment slack.Attachment) slack.Attachment {
	preText := fmt.Sprintf("New Rating for %s!", entry.Location)
	text := fmt.Sprintf("Timestamp: %s\nUser: %s\nRating: %s\nComment: %s", entry.Time.Format("2006-01-02"), entry.User, entry.Rating, entry.Comment)
	attachment = slack.Attachment{
		Title:      entry.Location,
		Pretext:    preText,
		Text:       text,
		TitleLink:  entry.Link,
		MarkdownIn: []string{"text", "pretext"},
	}
	return attachment
}
