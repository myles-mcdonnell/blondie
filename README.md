# Blondie

A command that will poll a TCP socket until a connection is made or the timeout is hit.  Useful if you need to wait for a server to start before doign something else. [https://www.youtube.com/watch?v=uWhkbDMISl8](https://www.youtube.com/watch?v=uWhkbDMISl8)

Args:

* timeout:  timeout in seconds, defaults to 5
* address: host and port, e.g. localhost:8080
* poll-interval: poll interval in milliseconds, defaults to 250

exit code is 1 if the connection is not established otherwise 0.
 
