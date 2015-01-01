package main

import(
	"fmt"
	"testing"
	"net/http"
	"net/http/httptest"
	"time"
)

func BenchmarkCallback(b *testing.B) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, "OK")
	}))

	defer ts.Close()

	service := Service{}

	callbackUri, err := service.ParseRule("callbacks/go"); if err != nil {
		b.Fatal("Failure could not parse rule: ", err.Error())
	}

	service.AttachHttpRule(callbackUri, ts.URL)
	go service.Listen("0.0.0.0", 8657, callbackUri)

	time.Sleep(1*time.Second)

	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			nreq, err := http.NewRequest("GET", "http://127.0.0.1:8657/callbacks/go?test=ok", nil)

			if err != nil {
			    b.Fatal("Failure making new request: ", err.Error())
			}

			client := http.Client{}

			resp, err := client.Do(nreq)

			if err != nil {
			    b.Fatal("Failure doing request: ", err.Error())
			}

			if resp.Status != "200 OK" {
				b.Fatal("Expected 200 got: ", resp.Status)
			}
		}
	})
}