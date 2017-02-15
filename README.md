# Blondie

A command that polls a TCP socket until a connection is made or the timeout is hit.  Useful if you need to wait for a server to start before doing something else.

Why call it Blondie? This is the best I could come up with in the 5 seconds I allowed for naming it [https://www.youtube.com/watch?v=uWhkbDMISl8](https://www.youtube.com/watch?v=uWhkbDMISl8)

Example:

```
//This will poll localhost:8080 every 500 milliseconds until a connection is made and it exits with code 0 and outputs 'Connection OK' or times out and exits code 2 with error message text output
blondie -address=localhost:8080 -timeout=20 -poll-interval=500 -exit-code-on-connect=0 exit-code-on-fail=2 -quiet-mode=false
```

Args:

* timeout:  timeout in seconds, defaults to 5
* address: host and port, e.g. localhost:8080
* poll-interval: poll interval in milliseconds, defaults to 250
* exit-code-on-connect: defaults to 0
* exit-code-on-fail: defaults to 1
* quiet-mode: suppress all output, default true

 
