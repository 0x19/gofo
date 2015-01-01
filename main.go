package main

import "flag"

var (
	host = flag.String("host", "0.0.0.0", "Host to listen on")
	port = flag.Int("port", 8070, "Port to listen on")
	rule = flag.String("in", "", "Rule that will be registered as server uri on which we will listen for requests")
	urls = flag.String("out", "", "List of URLS separated by comma. Same will receive http request received at rule")
)

func main() {
	flag.Parse()

	service := Service{}

	service.AttachHttpRule(*rule, *urls)
	service.Listen(*host, *port, *rule)
}