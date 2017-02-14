# Blondie

A command that will pool a TCP socket until a connection is made or the timeout expires.  Useful if you need to wait for a server to start before. [https://www.youtube.com/watch?v=uWhkbDMISl8](https://www.youtube.com/watch?v=uWhkbDMISl8)

Args:

* timeout:  timeout in seconds, defaults to 5
* address: host and port, e.g. localhost:8080
* poll-interval: poll interval in milliseconds, defaults to 250

exit code is 1 is the connection is not established otherwise 0
 
