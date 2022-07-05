package procrastiproxy

import (
	"errors"
	"flag"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"sort"
	"strings"
	"sync"
	"time"

	"github.com/hashicorp/go-multierror"

	"github.com/sirupsen/logrus"
	log "github.com/sirupsen/logrus"
)

// DefaultNow is the default implementation of Procrastiproxy's Now function,
// as we want procrastiproxy to return the actual time during normal operations. In
// testing, we override this method with static values (e.g., 9:00PM or 3:23 AM) in test
// cases to simulate different wall-times for verifiying procrastiproxy's behavior
var DefaultNow = time.Now

var logLevel string

type Procrastiproxy struct {
	Now  func() time.Time
	Port string
	List *List
	ProxyTimeSettings
}

type AdminCommand struct {
	Command string
	Host    string
}

type List struct {
	m       sync.Mutex
	members map[string]bool
}

type timeFlag struct {
	Name  string
	Value string
}

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

// RunCLI is the main entrypoint for the procrastiproxy package
func RunCLI() error {

	port := flag.String("port", "8000", "Port to listen on. Defaults to 8000")
	logLevel := flag.String("loglevel", "info", "Log level. Defaults to Info")
	blockList := flag.String("block", "", "Host to block. Defaults to none")
	blockStartTime := flag.String("block-start-time", defaultBlockStartTime, "Start of business hours. Defaults to 9:00AM")
	blockEndTime := flag.String("block-end-time", defaultBlockEndTime, "End of business hours. Defaults to 5:00PM")

	flag.Parse()

	level, err := log.ParseLevel(*logLevel)
	if err != nil {
		level = log.DebugLevel
	}
	log.SetLevel(level)

	p := NewProcrastiproxy()

	if parseErr := parseBlockListInput(blockList, p.GetList()); parseErr != nil {
		return parseErr
	}

	if parseErr := parseStartAndEndTimes(*blockStartTime, *blockEndTime); parseErr != nil {
		return parseErr
	}

	// Configure proxy time-based block settings
	p.ConfigureProxyTimeSettings(*blockStartTime, *blockEndTime)

	if *port == "" {
		return errors.New("You must supply a valid port via the --port flag")
	}
	if *blockList == "" {
		log.Debug("Proxy will allow all traffic, because you did not supply any sites to block via the --block flag")
	}

	p.SetPort(*port)

	RunServer(p)

	return nil
}

func makeProxyRequest(w http.ResponseWriter, r *http.Request) {
	// Perform the supplied request and return the response to caller
	res, err := http.Get(r.URL.String())
	if err != nil {
		log.Fatal(err)
	}
	body, err := ioutil.ReadAll(res.Body)
	// Return the request to caller by performing a very rudimentary clone of it, here
	w.Write(body)
}

func blockRequest(w http.ResponseWriter) {
	w.WriteHeader(http.StatusForbidden)
	w.Write([]byte("Forbidden"))
}

func (p *Procrastiproxy) proxyHandler(w http.ResponseWriter, r *http.Request) {
	log.WithFields(logrus.Fields{
		"Host": r.URL.Host,
		"Path": r.URL.Path,
	}).Debug("Proxy handler received request")

	log.WithFields(logrus.Fields{
		"blocked sites": p.GetList().All(),
	}).Debug("Blocked site hosts")

	makeProxyRequest(w, r)
}

func parseCommandFromPath(path string) (*AdminCommand, error) {
	aCmd := &AdminCommand{}
	pathElem := strings.Split(path, "/")
	if len(pathElem) < 4 {
		return aCmd, errors.New(fmt.Sprintf("Received malformed request path: %s\n", path))
	}
	if pathElem[2] == "block" || pathElem[2] == "unblock" {
		aCmd.Command = pathElem[2]
	}
	url, parseErr := url.Parse(pathElem[3])
	log.Debugf("Parsed URL: %s\n", url.String())
	if parseErr != nil {
		return aCmd, parseErr
	}
	aCmd.Host = url.String()
	return aCmd, nil
}

func (p *Procrastiproxy) timeAwareHandler(w http.ResponseWriter, r *http.Request) {
	if p.WithinBlockWindow(p.Now()) {
		log.Debug("Request made within block time window. Examining if host permitted..")
		p.blockListAwareHandler(w, r)
		return
	}
	log.Debug("Request made outside of configured block time window. Passing through...")
	p.proxyHandler(w, r)
}

func (p *Procrastiproxy) blockListAwareHandler(w http.ResponseWriter, r *http.Request) {
	host := sanitizeHost(r.URL.Host)
	if hostIsOnBlockList(host, p.GetList()) {
		log.Debugf("Blocking request to host: %s. User explicitly blocked and present time is within configured proxy block window", host)
		blockRequest(w)
		return
	}
	makeProxyRequest(w, r)
}

func (p *Procrastiproxy) adminHandler(w http.ResponseWriter, r *http.Request) {
	log.WithFields(logrus.Fields{
		"Path": r.URL.Path,
	}).Debug("Admin handler received request")

	adminCmd, err := parseCommandFromPath(r.URL.Path)
	if err != nil {
		log.Println(err)
	}

	var respMsg string
	list := p.GetList()

	if adminCmd.Command == "block" {
		list.Add(adminCmd.Host)
		respMsg = fmt.Sprintf("Successfully added: %s to the block list\n", adminCmd.Host)
	}
	if adminCmd.Command == "unblock" {
		list.Remove(adminCmd.Host)
		respMsg = fmt.Sprintf("Successfully removed: %s from the block list\n", adminCmd.Host)
	}
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(respMsg))
}

func NewList() *List {
	return &List{
		members: make(map[string]bool),
	}
}

// Clear resets the list, deleting all members
func (l *List) Clear() {
	defer l.m.Unlock()
	l.m.Lock()
	l.members = make(map[string]bool)
}

// All returns every member of the list
func (l *List) All() []string {
	l.m.Lock()
	defer l.m.Unlock()
	var members []string
	for k := range l.members {
		members = append(members, k)
	}
	return members
}

// Add appends an item to the list
func (l *List) Add(item string) {
	l.m.Lock()
	defer l.m.Unlock()
	l.members[item] = true
}

// Remove deletes an item from the list
func (l *List) Remove(item string) {
	l.m.Lock()
	defer l.m.Unlock()
	delete(l.members, item)
}

// Contains returns true if the supplied item is a member of the list
func (l *List) Contains(item string) bool {
	l.m.Lock()
	defer l.m.Unlock()
	return l.members[item]
}

// Length returns the number of members in the list
func (l *List) Length() int {
	l.m.Lock()
	defer l.m.Unlock()
	return len(l.members)
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

func sanitizeHost(host string) string {
	return strings.ToLower(strings.TrimSpace(strings.Replace(host, "\n", "", -1)))
}

func hostIsOnBlockList(host string, list *List) bool {
	host = sanitizeHost(host)
	return list.Contains(host)
}

func RunServer(p *Procrastiproxy) {

	log.WithFields(logrus.Fields{
		"Port":                    p.GetPort(),
		"Address":                 fmt.Sprintf("http://127.0.0.1:%s", p.GetPort()),
		"Number of sites blocked": p.GetList().Length(),
		"Log Level":               log.GetLevel().String(),
	}).Info("Procrastiproxy running...")

	http.HandleFunc("/", p.timeAwareHandler)
	http.HandleFunc("/admin/", p.adminHandler)

	log.Fatal(http.ListenAndServe(":"+p.GetPort(), nil))
}

func parseStartAndEndTimes(blockTimeStart, blockTimeEnd string) error {

	var result *multierror.Error

	timeFlags := []timeFlag{
		{
			Name:  "block-time-start",
			Value: blockTimeStart,
		},
		{
			Name:  "block-time-end",
			Value: blockTimeEnd,
		},
	}

	for _, timeFlag := range timeFlags {
		if _, parseErr := time.Parse(time.Kitchen, timeFlag.Value); parseErr != nil {
			result = multierror.Append(result, InvalidTimeFormatError{FlagName: timeFlag.Name, Value: timeFlag.Value, Underlying: parseErr})
		}
	}
	return result.ErrorOrNil()
}

// Build the fast, in-memory list of blocked hosts from the configured values
func AddHostToBlockList(list *List, hosts ...string) {
	for _, host := range hosts {
		list.Add(host)
	}
}

func SlicesAreEqual(a, b []string) bool {
	if len(a) != len(b) {
		return false
	}

	z := sort.StringSlice(a)
	y := sort.StringSlice(b)

	z.Sort()
	y.Sort()

	for i := range z {
		if z[i] != y[i] {
			return false
		}
	}

	return true
}

func validateBlockListInput(blockListMembers []string) error {
	if len(blockListMembers) < 1 {
		return EmptyBlockListError{}
	}
	return nil
}

func parseBlockListInput(blockList *string, list *List) error {
	var blockListMembers []string
	var blockListString = *blockList
	if blockListString != "" {
		blockListMembers = strings.Split(blockListString, ",")
	}

	validationErr := validateBlockListInput(blockListMembers)
	if validationErr != nil {
		return validationErr
	}

	for _, host := range blockListMembers {
		AddHostToBlockList(list, host)
	}

	return nil
}
