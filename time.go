package procrastiproxy

import (
	"time"

	"github.com/sirupsen/logrus"
	log "github.com/sirupsen/logrus"
)

var (
	proxyTimeSettings ProxyTimeSettings

	defaultBlockStartTime = "9:00AM"
	defaultBlockEndTime   = "5:00PM"
	defaultLayout         = "9:00AM"
)

type ProxyTimeSettings struct {
	Timezone       string
	BlockStartTime string
	BlockEndTime   string
	DefaultLayout  string
}

func ConfigureProxyTimeSettings(bts, bet string) {

	pts := ProxyTimeSettings{}
	if bts != "" {
		pts.BlockStartTime = bts
	} else {
		pts.BlockStartTime = defaultBlockStartTime
	}
	if bet != "" {
		pts.BlockEndTime = bet
	} else {
		pts.BlockEndTime = defaultBlockEndTime
	}
	pts.DefaultLayout = defaultLayout

	// Configure specified time zone
	log.WithFields(logrus.Fields{
		"Start block time": pts.BlockStartTime,
		"End block time":   pts.BlockEndTime,
	}).Infof("Successfully configured time settings")

	proxyTimeSettings = pts
}

func GetProxyTimeSettings() ProxyTimeSettings {
	if proxyTimeSettings == (ProxyTimeSettings{}) {
		// we haven't configured the settings and set the variable yet
		ConfigureProxyTimeSettings(defaultBlockStartTime, defaultBlockEndTime)
		return proxyTimeSettings
	}
	return proxyTimeSettings
}

func stringToTime(str string) time.Time {
	tm, err := time.Parse(time.Kitchen, str)
	if err != nil {
		log.Info("Failed to decode time:", err)
	}
	return tm
}

func WithinBlockWindow(check time.Time, pts ProxyTimeSettings) bool {
	startTimeString := pts.BlockStartTime
	endTimeString := pts.BlockEndTime

	timeNowString := check.Format(time.Kitchen)
	timeNow := stringToTime(timeNowString)

	start := stringToTime(startTimeString)
	end := stringToTime(endTimeString)

	if timeNow.Before(start) || timeNow.After(end) {
		return false
	}

	if timeNow.Before(end) {
		return true
	}

	return true
}
