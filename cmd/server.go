package cmd

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/sirupsen/logrus"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func hostIsBlocked(host string) bool {
	host = strings.ToLower(strings.TrimSpace(host))
	fmt.Printf("hostIsBlocked examining: %s\n", host)
	fmt.Printf("hostIsBlocked: %+v\n", viper.GetStringSlice("blocked_hosts"))
	for _, blockedHost := range viper.GetStringSlice("blocked_hosts") {
		fmt.Printf("blockedHost: %s and host: %s\n", blockedHost, host)
		if blockedHost == host {
			return true
		}
	}
	return false
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
		w.Write([]byte("Allowing: " + r.URL.Host))
	}
}

func RunServer(cmd *cobra.Command, args []string) {
	log.WithFields(logrus.Fields{
		"Port": args[0],
	}).Info("Proxy listening...")

	http.HandleFunc("/", proxyHandler)

	log.Fatal(http.ListenAndServe(":"+args[0], nil))
}
