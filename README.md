Go HTTP Fan Out
====
HTTP callback forwarder / fan-out. What it does is very simple. It listens any incoming HTTP packets on one uri and forwards what it gets to attached recipients.

This is quite useful if you need to listen for callbacks from external service which does not support notifying your application on multiple urls.
Just pass to external service generated url and bind your multiple urls. Once external service reaches generated url, attached urls will be instantly notified.


### Running behind NAT?

Well that's easy today :) You can always use [Ngrok](https://ngrok.com/)! It's free can very efficiently tunnel out your local machine. As well written in [Go](https://golang.org/)

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