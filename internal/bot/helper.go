package bot

import (
	"errors"
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
