package main

import (
	"flag"
	"github.com/op/go-logging"
)

var (
	log = logging.MustGetLogger("gofo")

	host = flag.String("host", "0.0.0.0", "Host to listen on")
	port = flag.Int("port", 8070, "Port to listen on")
	services = flag.String("services", "http", "List of fanout services that should be started up")
	rule = flag.String("in", "", "Rule that will be registered as server uri on which we will listen for requests")
	urls = flag.String("out", "", "List of URLS separated by comma. Same will receive http request received at rule")
)

func main() {
	flag.Parse()

	logging.MustStringFormatter("%{color}%{shortfunc} ▶ %{level:.8s} %{color:reset} %{message}")

	service := Service{}

	callbackUri, err := service.ParseRule(*rule); if err != nil {
		log.Error("Error while parsing rule: %s", err.Error())
		return
	}

	// I guess if url is smaller than 10 same means that it's smaller than http://t.t
	if len(*urls) < 10 {
		log.Error("You need to provide valid urls into out in order to start service!")
		return
	}

	service.AttachHttpRule(callbackUri, *urls)
	service.Listen(*host, *port, callbackUri)
}