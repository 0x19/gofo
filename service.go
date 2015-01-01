package main

import (
	"io"
	"bytes"
	"fmt"
	"net/http"
	"net/url"
	"strings"
	"regexp"
	"errors"
)

type Service struct {
	Rules []string
}

func(s *Service) processForwarder(furl *url.URL, req *http.Request) {
	log.Debug("Processing forwarder: %s", furl)

	fullUrl := furl.String()
	var body io.Reader

	if req.URL.RawQuery != "" && req.Method == "GET" {
		fullUrl = fmt.Sprintf("%s?%s", fullUrl, req.URL.RawQuery)
	}

	if req.Method != "GET" {
		body = bytes.NewBufferString(req.Form.Encode())
	}

	nreq, err := http.NewRequest(req.Method, fullUrl, body)

	if req.Method != "GET" {
		nreq.Header.Set("Content-Type", "application/x-www-form-urlencoded;")
	}

	client := http.Client{}

	if err != nil {
	    panic(err)
	}

	resp, err := client.Do(nreq)

	if err != nil {
	    panic(err)
	}

	log.Info("Forwarder response status: %q", resp.Status)
    log.Info("Forwarder response headers: %s", resp.Header)
}

func(s *Service) isURL(url string) bool {
	if url == "" || len(url) >= 2083 || len(url) < 10 {
		return false
	}

	rxURL := regexp.MustCompile(`^((http|https):\/\/)?(\S+(:\S*)?@)?((([1-9]\d?|1\d\d|2[01]\d|22[0-3])(\.(1?\d{1,2}|2[0-4]\d|25[0-5])){2}(?:\.([0-9]\d?|1\d\d|2[0-4]\d|25[0-4]))|((www\.)?)?(([a-z\x{00a1}-\x{ffff}0-9]+-?-?_?)*[a-z\x{00a1}-\x{ffff}0-9]+)(?:\.([a-z\x{00a1}-\x{ffff}]{2,}))?)|localhost)(:(\d{1,5}))?((\/|\?|#)[^\s]*)?$`)
	return rxURL.MatchString(url)
}

func(s *Service) prepareForwarders(fwds string) []*url.URL {
	forwarders := make([]*url.URL, 0)

	for _, forwarder := range strings.Split(fwds, ",") {
		forwarder = strings.Trim(forwarder, " ")
		log.Debug("Validating attached forwarder: %s", forwarder)

		if !s.isURL(forwarder) {
			panic(fmt.Errorf("Hola, passed forwarder `%s` is not valid url! Only valid urls can be attached as forwarders.", forwarder))
		}

		urlfwd, err := url.Parse(forwarder)

		if err != nil {
			panic(err)
		}

		forwarders = append(forwarders, urlfwd)
	}

	return forwarders
}

// Will parse rule and assign it to the struct. If however is not valid will raise error saying same
func(s *Service) ParseRule(rule string) (string, error) {
	if !strings.HasPrefix(rule, "/") {
		rule = fmt.Sprintf("/%s", rule)
	}

	rxAn := regexp.MustCompile("^[a-zA-Z0-9-_/.]+$") // AlphaNumeric + _ and -

	if !rxAn.MatchString(rule) {
		return rule, errors.New("Rule contains illegal characters. Accepted: a-z A-Z 0-9 - _ / and .")
	}

	return rule, nil
}

func(s *Service) AttachHttpRule(rule string, fwds string) {

	forwarders := s.prepareForwarders(fwds)
	log.Debug("Attaching new rule: %s -> %s", rule, forwarders)

	for _, r := range s.Rules {
		if r == rule {
			log.Debug("Requested rule is already attached. You cannot re-attach already attached rule: %s", rule)
			return
		}
	}

	s.Rules = append(s.Rules, rule)

	http.HandleFunc(rule, func(w http.ResponseWriter, r *http.Request) {
		if err := r.ParseForm(); err != nil {
			fmt.Fprintf(w, "error parsing form.")
			return
		}

		fmt.Fprintf(w, "OK")

		for _, forwarder := range forwarders {
			go s.processForwarder(forwarder, r)
		}
	})
}

func(s *Service) Listen(host string, port int, rule string) {
	log.Notice("Listening for new incoming connections %s:%d", host, port)
	log.Notice("You can pass following URL to external service: http://%s:%d%s", host, port, rule)
	http.ListenAndServe(fmt.Sprintf("%s:%d", host, port), nil)
}