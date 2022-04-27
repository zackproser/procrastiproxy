package procrastiproxy

import (
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"

	"github.com/sirupsen/logrus"
	log "github.com/sirupsen/logrus"
)

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

func proxyHandler(w http.ResponseWriter, r *http.Request) {
	log.WithFields(logrus.Fields{
		"Host": r.URL.Host,
		"Path": r.URL.Path,
	}).Debug("Proxy handler received request")

	log.WithFields(logrus.Fields{
		"blocked sites": GetList().All(),
	}).Debug("Blocked site hosts")

	makeProxyRequest(w, r)
}

type AdminCommand struct {
	Command string
	Host    string
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
	log.Printf("Parsed URL: %s\n", url.String())
	if parseErr != nil {
		return aCmd, parseErr
	}
	aCmd.Host = url.String()
	return aCmd, nil
}

func (p *Procrastiproxy) timeAwareHandler(w http.ResponseWriter, r *http.Request) {
	if p.WithinBlockWindow() {
		log.Info("Request made within block time window. Examining if host permitted..")
		p.blockListAwareHandler(w, r)
		return
	}
	log.Info("Request made outside of configured block time window. Passing through...")
	proxyHandler(w, r)
}

func (p *Procrastiproxy) blockListAwareHandler(w http.ResponseWriter, r *http.Request) {
	host := sanitizeHost(r.URL.Host)
	if hostIsBlocked(host) {
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
	list := GetList()
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
