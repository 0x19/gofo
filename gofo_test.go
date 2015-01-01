package main

import(
	"fmt"
	"testing"
	"net/http"
	"net/http/httptest"
)


func TestRule(t *testing.T) {
	service := Service{}

	callbackUri, err := service.ParseRule("callbacks/go"); if err != nil {
		t.Fatal("Failure could not parse rule: ", err.Error())
	}

	// Testing if / will be applied
	if callbackUri != "/callbacks/go" {
		t.Fatal("Parsing failed due to / not auto attached to the rule")
	}

	_, err = service.ParseRule("callbacks/go≤"); if err == nil {
		t.Fatal("Failed because rule contains invalid character")
	}

	_, err = service.ParseRule("callbacks/go-this_should_be_good-a./"); if err != nil {
		t.Fatal("Failed because of: ", err.Error())
	}
}

func TestCallback(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, "OK")
	}))

	defer ts.Close()

	service := Service{}

	callbackUri, err := service.ParseRule("callbacks/go"); if err != nil {
		t.Fatal("Failure could not parse rule: ", err.Error())
	}

	service.AttachHttpRule(callbackUri, ts.URL)
	go service.Listen("0.0.0.0", 8657, callbackUri)

	nreq, err := http.NewRequest("GET", "http://127.0.0.1:8657/callbacks/go?test=ok", nil)

	if err != nil {
	    t.Fatal("Failure making new request: ", err.Error())
	}

	client := http.Client{}

	resp, err := client.Do(nreq)

	if err != nil {
	    t.Fatal("Failure doing request: ", err.Error())
	}

	if resp.Status != "200 OK" {
		t.Fatal("Expected 200 got: ", resp.Status)
	}
}