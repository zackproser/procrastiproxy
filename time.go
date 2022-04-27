package procrastiproxy

import (
	"fmt"
	"os"
	"time"

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

func (p *Procrastiproxy) ConfigureProxyTimeSettings(bts, bet string) {

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

	p.ProxyTimeSettings = pts
}

func (p *Procrastiproxy) GetProxyTimeSettings() ProxyTimeSettings {
	if p.ProxyTimeSettings == (ProxyTimeSettings{}) {
		// we haven't configured the settings and set the variable yet
		p.ConfigureProxyTimeSettings(defaultBlockStartTime, defaultBlockEndTime)
		return p.ProxyTimeSettings
	}
	return p.ProxyTimeSettings
}

func stringToTime(str string) time.Time {
	if str == "" {
		fmt.Println("its empty")
		os.Exit(1)
	}
	tm, err := time.Parse(time.Kitchen, str)
	if err != nil {
		log.Infof("Failed to decode time: %s - error: %v\n", str, err)
	}
	return tm
}

var now = time.Now

func (p *Procrastiproxy) WithinBlockWindow() bool {

	check := p.Now()
	//fmt.Printf("check: %v\n", check)
	pts := p.GetProxyTimeSettings()
	//fmt.Printf("pts: %v\n", pts)

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
