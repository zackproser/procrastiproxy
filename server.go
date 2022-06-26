package procrastiproxy

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/sirupsen/logrus"
	log "github.com/sirupsen/logrus"
)

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
