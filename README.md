# Blondie

[![Go Report Card](https://goreportcard.com/badge/github.com/myles-mcdonnell/blondie)](https://goreportcard.com/report/github.com/myles-mcdonnell/blondie)


A command that polls one or more TCP sockets until a connection is made or the timeout is hit.  Useful if you need to wait for a server or group of servers to start before doing something else.

Why call it Blondie? This is the best I could come up with in the 5 seconds I allowed for naming it [https://www.youtube.com/watch?v=uWhkbDMISl8](https://www.youtube.com/watch?v=uWhkbDMISl8)

Example:

```
//This will poll localhost:8080 and localhost:80 every 500 milliseconds until a connection is made and it exits with code 0 and outputs 'Connection OK' or the timeout is reached (20 and 5 seconds respectively) and exits code 2 with error message text output
blondie -targets=localhost:8080:20,localhost:80:5 -poll-interval=500 -exit-code-on-connect=0 exit-code-on-fail=2 -quiet-mode=false
```

Args:

* targets: comma separated host, port and timeoutSeconds, e.g. localhost:8080:20,localhost:80:5
* globalTimeout: global timeout in seconds
* poll-interval: poll interval in milliseconds, defaults to 250
* exit-code-on-connect: defaults to 0
* exit-code-on-fail: defaults to 1
* quiet-mode: suppress all output, default true

 
