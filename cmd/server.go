package cmd

import (
	"net/http"
	"strings"

	"github.com/sirupsen/logrus"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

func hostIsBlocked(host string) bool {
	host = strings.ToLower(strings.TrimSpace(host))
	blockList := GetList()
	return blockList.Contains(host)
}

func RunServer(cmd *cobra.Command, args []string) {
	log.WithFields(logrus.Fields{
		"Port": args[0],
	}).Info("Proxy listening...")

	http.HandleFunc("/", proxyHandler)
	http.HandleFunc("/admin/", adminHandler)

	log.Fatal(http.ListenAndServe(":"+args[0], nil))
}
