package procrastiproxy

import (
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

// stringToTime accepts a string representation of a timestamp and attempts to convert it to
// a time in the "Kitchen" format, e.g., 3:04PM
func stringToTime(str string) time.Time {
	tm, err := time.Parse(time.Kitchen, str)
	if err != nil {
		log.Debugf("Failed to decode time: %s - error: %v\n", str, err)
	}
	return tm
}

func (p *Procrastiproxy) WithinBlockWindow(now time.Time) bool {

	pts := p.GetProxyTimeSettings()

	startTimeString := pts.BlockStartTime
	endTimeString := pts.BlockEndTime

	start := stringToTime(startTimeString)
	end := stringToTime(endTimeString)

	// Create an equivalent unix epoch timestamp, but use now's hour, minutes and seconds
	checkTime := time.Date(int(0000), time.January, int(1), now.Hour(), now.Minute(), now.Second(), now.Nanosecond(), time.UTC)

	log.Debugf("startTime: %v endTime: %v currentTime: %v", start, end, checkTime)

	if checkTime.Before(start) {
		return false
	}

	if checkTime.Before(end) {
		return true
	}

	return false

}
