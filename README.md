[![License](http://img.shields.io/badge/license-MIT-blue.svg?style=flat)]()
[![Build Status](https://travis-ci.org/0x19/gofo.svg)](https://travis-ci.org/0x19/gofo)
[![Go 1.3 Ready](https://img.shields.io/badge/Go%201.3-Ready-green.svg?style=flat)]()
[![Go 1.4 Ready](https://img.shields.io/badge/Go%201.4-Ready-green.svg?style=flat)]()

HTTP Fan Out Server written in Go
====
What it does is very simple. It listens for incoming HTTP requests on registered uri (-in) and forwards what it gets to registered listeners (-out).

This is quite useful if you need to listen for callbacks from external service if same one does not support notifying your application on multiple urls.

Whenever you start service it will give you full url that you can just copy and paste to external service.

```shell
2015/01/01 17:35:01 Listen ▶ NOTICE  You can pass following URL to external service: http://127.0.0.1:9991/callbacks/sms
```


### Running behind NAT?

Well that's easy today :) You can always use [Ngrok](https://ngrok.com/)! It's free and can for sure very efficiently tunnel out your local machine. As well written in [Go](https://golang.org/)


### Issues?

Please reach me over nevio.vesic@gmail.com or submit new Issue. I'd prefer tho if you would submit issue.

### Example


```shell
./gofo -host 127.0.0.1 -port 9991 -in callbacks/sms -out "http://webhookr.com/one, http://webhookr.com/one/two, http://webhookr.com/one/three"

2015/01/01 17:35:01 prepareForwarders ▶ DEBUG  Validating attached forwarder: http://webhookr.com/one
2015/01/01 17:35:01 prepareForwarders ▶ DEBUG  Validating attached forwarder: http://webhookr.com/one/two
2015/01/01 17:35:01 prepareForwarders ▶ DEBUG  Validating attached forwarder: http://webhookr.com/one/three
2015/01/01 17:35:01 AttachHttpRule ▶ DEBUG  Attaching new rule: /callbacks/sms -> [http://webhookr.com/one http://webhookr.com/one/two http://webhookr.com/one/three]
2015/01/01 17:35:01 Listen ▶ NOTICE  Listening for new incoming connections 127.0.0.1:9991
2015/01/01 17:35:01 Listen ▶ NOTICE  You can pass following URL to external service: http://127.0.0.1:9991/callbacks/sms
```
