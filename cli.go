/*
Copyright © 2022 Zack Proser zackproser@gmail.com

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/
package procrastiproxy

import (
	"errors"
	"flag"

	log "github.com/sirupsen/logrus"
)

var (
	logLevel string
)

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

	parseBlockListInput(blockList)

	p := NewProcrastiproxy()
	// Configure proxy time-based block settings
	p.ConfigureProxyTimeSettings(*blockStartTime, *blockEndTime)

	if *port == "" {
		return errors.New("You must supply a valid port via the --port flag")
	}
	if *blockList == "" {
		log.Info("Proxy will allow all traffic, because you did not supply any sites to block via the --block flag")
	}
	args := []string{*port, *blockList}
	RunServer(args)
	return nil
}
