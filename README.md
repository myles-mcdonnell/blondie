## Blondie

[![Go Report Card](https://goreportcard.com/badge/github.com/myles-mcdonnell/blondie)](https://goreportcard.com/report/github.com/myles-mcdonnell/blondie)

Blondie is a CLI tool and package that tries to connect to HTTP and TCP endpoints up to a configured timeout.

The success criteria in the case of a TCP socket is a successful connection.  HTTP checks can be configured to accept any or one from a list of status codes as successful.


### Install
The latest binaries for all supported operating systems are [here](https://github.com/myles-mcdonnell/blondie/releases)

You may choose to execute the following script which will download the latest binary to `/usr/local/bin/blondie`:

```
sudo curl https://raw.githubusercontent.com/myles-mcdonnell/blondie/master/install.sh | sudo sh
```


If you have Go tool installed you may also run:
```
go get -u github.com/myles-mcdonnell/blondie/...
```


### Usage
```sh
blondie -help
```

Example CLI usage:

```
blondie -targets=tcp:icanhazip.com:80:2000,http:icanhazip.com:80:5000,https:github.com:443:4000::200_204 -poll-interval=500 -exit-code-on-connect=0 exit-code-on-fail=2 -quiet-mode=false
```

Args:

* targets: comma separated host, port, timeoutMilliseconds, path and successHttpCode, e.g. tcp:localhost:8080:20000,http:localhost:80:5000::200,http:localhost:80:5000:health:200_204,http:localhost:80:5 (path and successHttpCode optional for http protocol)
* poll-interval: poll interval in milliseconds, defaults to 250
* exit-code-on-connect: defaults to 0
* exit-code-on-fail: defaults to 1
* quiet-mode: suppress all output, default true

### Naming
Why call it Blondie? [https://www.youtube.com/watch?v=uWhkbDMISl8](https://www.youtube.com/watch?v=uWhkbDMISl8)

 
