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

func(s *Service) handleUrl(furl *url.URL, req *http.Request) {
	log.Debug("About to handle url: %s", furl)

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

	log.Debug("Forwarder response status: %q", resp.Status)
    log.Debug("Forwarder response headers: %s", resp.Header)
}

// Will check whenever given url is in fact url
func(s *Service) isURL(url string) bool {
	if url == "" || len(url) >= 2083 || len(url) < 10 {
		return false
	}

	rxURL := regexp.MustCompile(`^((http|https):\/\/)?(\S+(:\S*)?@)?((([1-9]\d?|1\d\d|2[01]\d|22[0-3])(\.(1?\d{1,2}|2[0-4]\d|25[0-5])){2}(?:\.([0-9]\d?|1\d\d|2[0-4]\d|25[0-4]))|((www\.)?)?(([a-z\x{00a1}-\x{ffff}0-9]+-?-?_?)*[a-z\x{00a1}-\x{ffff}0-9]+)(?:\.([a-z\x{00a1}-\x{ffff}]{2,}))?)|localhost)(:(\d{1,5}))?((\/|\?|#)[^\s]*)?$`)
	return rxURL.MatchString(url)
}

// Will parse url and append to the slice
func(s *Service) parseUrls(out string) []*url.URL {
	urls := make([]*url.URL, 0)

	for _, fw := range strings.Split(out, ",") {
		fw = strings.Trim(fw, " ")
		log.Debug("Validating attached forwarder: %s", fw)

		if !s.isURL(fw) {
			panic(fmt.Errorf("Hola, passed forwarder `%s` is not valid url! Only valid urls can be attached.", fw))
		}

		pfw, err := url.Parse(fw)

		if err != nil {
			panic(err)
		}

		urls = append(urls, pfw)
	}

	return urls
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

// Will attach rule and its urls and prepare for listening
func(s *Service) HandleFanRequest(rule string, out string) {

	urls := s.parseUrls(out)
	log.Debug("Attaching new rule: %s -> %s", rule, urls)

	for _, r := range s.Rules {
		if r == rule {
			panic(fmt.Sprintf("Requested rule is already attached. You cannot re-attach already attached rule: %s", rule))
		}
	}

	s.Rules = append(s.Rules, rule)

	http.HandleFunc(rule, func(w http.ResponseWriter, r *http.Request) {
		if err := r.ParseForm(); err != nil {
			fmt.Fprintf(w, "error parsing form.")
			return
		}

		fmt.Fprintf(w, "OK")

		for _, url := range urls {
			go s.handleUrl(url, r)
		}
	})
}

// Will nuke mother fu**er!
func(s *Service) Listen(host string, port int, rule string) {
	log.Notice("Listening for new incoming connections %s:%d", host, port)
	log.Notice("You can pass following URL to external service: http://%s:%d%s", host, port, rule)
	http.ListenAndServe(fmt.Sprintf("%s:%d", host, port), nil)
}