package procrastiproxy

import (
	"net/http"
	"strings"

	"github.com/sirupsen/logrus"
	log "github.com/sirupsen/logrus"
)

func sanitizeHost(host string) string {
	return strings.ToLower(strings.TrimSpace(strings.Replace(host, "\n", "", -1)))
}

func hostIsBlocked(host string) bool {
	host = sanitizeHost(host)
	blockList := GetList()
	return blockList.Contains(host)
}

func RunServer(args []string) {
	port := args[0]
	log.WithFields(logrus.Fields{
		"Port": port,
	}).Info("Starting server on port...")

	http.HandleFunc("/", timeAwareHandler)
	http.HandleFunc("/admin/", adminHandler)

	log.Fatal(http.ListenAndServe(":"+port, nil))
}
