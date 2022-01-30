package cmd

import (
	"time"

	"github.com/sirupsen/logrus"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

var proxyTimeSettings ProxyTimeSettings

type ProxyTimeSettings struct {
	Timezone       string
	BlockStartTime string
	BlockEndTime   string
	DefaultLayout  string
}

func ConfigureProxyTimeSettings() {
	var (
		defaultTimezone       = "America/New_York"
		defaultBlockStartTime = "9:00AM"
		defaultBlockEndTime   = "5:00PM"
		defaultLayout         = "9:00AM"
	)

	pts := ProxyTimeSettings{}
	if tz := viper.GetString("timezone"); tz != "" {
		pts.Timezone = tz
	} else {
		pts.Timezone = defaultTimezone
	}
	if bts := viper.GetString("block_start_time"); bts != "" {
		pts.BlockStartTime = bts
	} else {
		pts.BlockStartTime = defaultBlockStartTime
	}
	if bet := viper.GetString("block_end_time"); bet != "" {
		pts.BlockEndTime = bet
	} else {
		pts.BlockEndTime = defaultBlockEndTime
	}
	pts.DefaultLayout = defaultLayout

	// Configure specified time zone
	log.WithFields(logrus.Fields{
		"Timezone":         pts.Timezone,
		"Start block time": pts.BlockStartTime,
		"End block time":   pts.BlockEndTime,
	}).Infof("Successfully configured time settings")

	proxyTimeSettings = pts
}

func GetProxyTimeSettings() ProxyTimeSettings {
	if proxyTimeSettings == (ProxyTimeSettings{}) {
		// we haven't configured the settings and set the variable yet
		ConfigureProxyTimeSettings()
		return proxyTimeSettings
	}
	return proxyTimeSettings
}

func stringToTime(str string) time.Time {
	tm, err := time.Parse(time.Kitchen, str)
	if err != nil {
		log.Info("Failed to decode time:", err)
	}
	log.Info("Time decoded:", tm)
	return tm
}

func WithinBlockWindow(check time.Time, pts ProxyTimeSettings) bool {
	startTimeString := pts.BlockStartTime
	endTimeString := pts.BlockEndTime

	timeNowString := check.Format(time.Kitchen)
	timeNow := stringToTime(timeNowString)

	start := stringToTime(startTimeString)
	end := stringToTime(endTimeString)

	if timeNow.Before(start) {
		return false
	}

	if timeNow.Before(end) {
		return true
	}

	return true
}
