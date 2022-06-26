package procrastiproxy

import (
	"fmt"
	"time"
)

// DefaultNow is the default implementation of Procrastiproxy's Now function,
// as we want procrastiproxy to return the actual time during normal operations. In
// testing, we override this method with static values (e.g., 9:00PM or 3:23 AM) in test
// cases to simulate different wall-times for verifiying procrastiproxy's behavior
var DefaultNow = time.Now

type Procrastiproxy struct {
	Now  func() time.Time
	Port string
	List *List
	ProxyTimeSettings
}

func NewProcrastiproxy() *Procrastiproxy {
	return &Procrastiproxy{
		Now:  DefaultNow,
		List: NewList(),
	}
}

func (p *Procrastiproxy) GetList() *List {
	return p.List
}

func (p *Procrastiproxy) SetPort(s string) {
	p.Port = s
}

func (p *Procrastiproxy) GetPort() string {
	return p.Port
}

// custom errors

type EmptyBlockListError struct{}

func (err EmptyBlockListError) Error() string {
	return fmt.Sprint("You must supply at least one valid HTTP host to procrastiproxy via the --block flag. Example: --block reddit.com")
}

type InvalidTimeFormatError struct {
	FlagName   string
	Value      string
	Underlying error
}

func (err InvalidTimeFormatError) Error() string {
	return fmt.Sprintf("Invalid time value {%s} passed with flag {%s}. Format must be time.Kitchen: e.g., 9:15AM. Parse error: %v", err.Value, err.FlagName, err.Underlying)
}
