package main

import "github.com/op/go-logging"

var log = logging.MustGetLogger("gomany")

func init() {
	logging.MustStringFormatter("%{color}%{shortfunc} ▶ %{level:.8s} %{color:reset} %{message}")
}