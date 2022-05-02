package procrastiproxy

import "time"

var procrastiproxy *Procrastiproxy

// DefaultNow is the default implementation of Procrastiproxy's Now function,
// as we want procrastiproxy to return the actual time during normal operations. In
// testing, we override this method with static values (e.g., 9:00PM or 3:23 AM) in test
// cases to simulate different wall-times for verifiying procrastiproxy's behavior
var DefaultNow = time.Now

type Procrastiproxy struct {
	Now func() time.Time
	ProxyTimeSettings
}

func NewProcrastiproxy() *Procrastiproxy {
	if procrastiproxy != nil {
		return procrastiproxy
	}
	procrastiproxy = &Procrastiproxy{
		Now: DefaultNow,
	}
	return procrastiproxy
}
