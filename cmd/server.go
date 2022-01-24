package cmd

import (
	"io/ioutil"
	"net/http"
	"strings"

	"github.com/sirupsen/logrus"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func hostIsBlocked(host string) bool {
	host = strings.ToLower(strings.TrimSpace(host))
	blockList := GetList()
	return blockList.Contains(host)
}

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

func RunServer(cmd *cobra.Command, args []string) {
	log.WithFields(logrus.Fields{
		"Port": args[0],
	}).Info("Proxy listening...")

	http.HandleFunc("/", proxyHandler)
	http.HandleFunc("/admin/", adminHandler)

	log.Fatal(http.ListenAndServe(":"+args[0], nil))
}
