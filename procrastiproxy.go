package procrastiproxy

import "time"

var procrastiproxy *Procrastiproxy

type Procrastiproxy struct {
	Now func() time.Time
	ProxyTimeSettings
}

func NewProcrastiproxy() *Procrastiproxy {
	if procrastiproxy != nil {
		return procrastiproxy
	}
	procrastiproxy = &Procrastiproxy{}
	return procrastiproxy
}
