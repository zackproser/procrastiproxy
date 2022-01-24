package cmd

import (
	"io/ioutil"
	"net/http"
	"strings"

	"github.com/sirupsen/logrus"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

func proxyHandler(w http.ResponseWriter, r *http.Request) {
	log.WithFields(logrus.Fields{
		"Host": r.URL.Host,
		"Path": r.URL.Path,
	}).Debug("Proxy handler received request")

	log.WithFields(logrus.Fields{
		"blocked sites": viper.Get("blocked_hosts"),
	}).Debug("Blocked site hosts")

	if hostIsBlocked(r.URL.Host) {
		w.WriteHeader(http.StatusForbidden)
		w.Write([]byte("Forbidden"))
	} else {
		// Perform the supplied request and return the response to caller
		res, err := http.Get(r.URL.String())
		if err != nil {
			log.Fatal(err)
		}
		body, err := ioutil.ReadAll(res.Body)
		// Return the request to caller by performing a very rudimentary clone of it, here
		w.WriteHeader(res.StatusCode)
		w.Write(body)
	}
}

func adminHandler(w http.ResponseWriter, r *http.Request) {
	log.WithFields(logrus.Fields{
		"Path": r.URL.Path,
	}).Debug("Admin handler received request")
	pathElem := strings.Split(r.URL.Path, "/")
	log.Printf("pathElem: %+v pathElem[3]: %+v\n", pathElem, pathElem[2])
	if len(pathElem) < 1 {
		log.Printf("Received malformed request path: %s\n", r.URL.Path)
	}
	if pathElem[2] == "block" || pathElem[2] == "unblock" {
		log.Printf("Received valid command: %s\n", pathElem[2])
	}
}
