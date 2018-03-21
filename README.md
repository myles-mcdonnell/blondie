# Blondie

[![Go Report Card](https://goreportcard.com/badge/github.com/myles-mcdonnell/blondie)](https://goreportcard.com/report/github.com/myles-mcdonnell/blondie)


## install
```sh
go get -u github.com/banaio/blondie/...
blondie -help
Usage of blondie:
  -exit-code-on-connect int
    	Exit code when connection is made
  -exit-code-on-fail int
    	Exit code when connection is not made (default 1)
  -poll-interval int
    	poll interval in milliseconds (default 250)
  -quiet-mode
    	suppress all output
  -targets string
    	comma separated protocol:address:port:timeoutSeconds:[path]:[successcode], e.g. [tcp|http|https]:localhost:8080:60:[path]:[success_error_code] timeout is optional in which case the global timeout is used, final two arguments for htyp only
```

A command that polls one or more TCP sockets/HTTP endpoints until a connection is made or the timeout is hit.  Useful if you need to wait for a server or group of servers to start before doing something else.

Why call it Blondie? This is the best I could come up with in the 5 seconds I allowed for naming it [https://www.youtube.com/watch?v=uWhkbDMISl8](https://www.youtube.com/watch?v=uWhkbDMISl8)

Example:

```
//This will poll tcp:localhost:8080 and http:localhost:80 every 500 milliseconds until a connection is made and it exits with code 0 and outputs 'Connection OK' or the timeout is reached (20 and 5 seconds respectively) and exits code 2 with error message text output
blondie -targets=tcp:localhost:8080:2000,http:localhost:80:5000 -poll-interval=500 -exit-code-on-connect=0 exit-code-on-fail=2 -quiet-mode=false
```

Args:

* targets: comma separated host, port, timeoutMilliseconds, path and successHttpCode, e.g. tcp:localhost:8080:20,http:localhost:80:5::200,http:localhost:80:5:health:200,http:localhost:80:5 (path and successHttpCode optional for http protocol)
* poll-interval: poll interval in milliseconds, defaults to 250
* exit-code-on-connect: defaults to 0
* exit-code-on-fail: defaults to 1
* quiet-mode: suppress all output, default true

 
