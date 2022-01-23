package cmd

import (
	"fmt"
	"net/http"

	"github.com/sirupsen/logrus"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func RunServer(cmd *cobra.Command, args []string) {
	log.WithFields(logrus.Fields{
		"Port": viper.Get("port"),
	}).Info("Proxy listening...")

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "Pong")
	})

	log.Fatal(http.ListenAndServe(":"+viper.GetString("port"), nil))
}
