package cmd

import (
	"fmt"
	"time"

	"github.com/sirupsen/logrus"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

var proxyTimeSettings *ProxyTimeSettings

type ProxyTimeSettings struct {
	Timezone       string
	BlockStartTime string
	BlockEndTime   string
}

func ConfigureProxyTimeSettings() {
	var (
		defaultTimezone       = "America/New_York"
		defaultBlockStartTime = "9AM"
		defaultBlockEndTime   = "5PM"
	)

	pts := &ProxyTimeSettings{}
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

	// Configure specified time zone
	log.WithFields(logrus.Fields{
		"Timezone":         pts.Timezone,
		"Start block time": pts.BlockStartTime,
		"End block time":   pts.BlockEndTime,
	}).Infof("Successfully configured time settings")

	proxyTimeSettings = pts
}

func GetProxyTimeSettings() *ProxyTimeSettings {
	if proxyTimeSettings == (&ProxyTimeSettings{}) {
		// we haven't configured the settings and set the variable yet
		ConfigureProxyTimeSettings()
		return proxyTimeSettings
	}
	return proxyTimeSettings
}

func WithinBlockWindow() bool {
	pts := GetProxyTimeSettings()
	location, err := time.LoadLocation(pts.Timezone)
	if err != nil {
		log.Infof("Error getting time location from: %s\n", pts.Timezone)
	}
	tm, parseErr := time.ParseInLocation("05:04:00PM", pts.BlockStartTime, location)
	if parseErr != nil {
		log.Infof("Error parsing time in location: %+v\n", parseErr)
	}
	t := time.Now().In(location)
	fmt.Printf("time settings: %+v\n", pts)
	fmt.Printf("time.Now().In(location): %v\n", t)
	fmt.Printf("t.After(tm) %v\n", t.After(tm))
	fmt.Printf("t.Before(tm) %v\n", t.Before(tm))
	return true
}
