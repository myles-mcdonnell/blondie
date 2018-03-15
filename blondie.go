package main

import (
	"flag"
	"fmt"
	"net"
	"os"
	"strconv"
	"strings"
	"sync"
	"time"
	"net/http"
)

type outputWriter struct {
	QuietMode bool
}

func main() {

	targets := flag.String("targets", "", "comma separated protocol:address:port:timeoutSeconds, e.g. [tcp|http|https]:localhost:8080:60:[path]:[success_error_code] timeout is optional in which case the global timeout is used, final two arguments for htyp only")
	globalTimeout := flag.Float64("globalTimeout", 5, "a global timeoput")
	pollinterval := flag.Int("poll-interval", 250, "poll interval in milliseconds")
	exitCodeOnConnectOk := flag.Int("exit-code-on-connect", 0, "Exit code when connection is made")
	exitCodeOnConnectFail := flag.Int("exit-code-on-fail", 1, "Exit code when connection is not made")
	quietMode := flag.Bool("quiet-mode", false, "suppress all output")
	flag.Parse()

	if len(*targets) < 1 {
		showUsage()
		return
	}

	var waitGroup = &sync.WaitGroup{}
	var writer = &outputWriter{QuietMode: *quietMode}
	var failure bool

	for _, target := range strings.Split(*targets, ",") {
		protoAddrPortAndTimeout := strings.Split(target, ":")

		var protocol = strings.ToLower(protoAddrPortAndTimeout[0])

		if protocol!="tcp"&&protocol!="http"&&protocol!="https" {
			fmt.Printf("Unrecognised protocol %s", protocol)
			os.Exit(*exitCodeOnConnectFail)
		}

		var address = protoAddrPortAndTimeout[1] + ":" + protoAddrPortAndTimeout[2]

		if address == "" {
			fmt.Println("Please supply targets, e.g. -targets=localhost:8080;60,google.com:80;30")
			os.Exit(*exitCodeOnConnectFail)
		}

		var path string
		var successcode int
		if protocol=="http"||protocol=="https" {
			path = 	protoAddrPortAndTimeout[4]
			var err error
			successcode, err = strconv.Atoi(protoAddrPortAndTimeout[5])
			if err!=nil {
				fmt.Println("Can not parse success code to int %s", protoAddrPortAndTimeout[5] )
				os.Exit(*exitCodeOnConnectFail)
			}
		}

		var timeout float64
		if len(protoAddrPortAndTimeout) < 4 {
			timeout = *globalTimeout
		} else {
			var err error
			to, err := strconv.Atoi(protoAddrPortAndTimeout[3])
			if err != nil {
				panic(err)
			}

			timeout = float64(to)
		}

		start := time.Now()

		waitGroup.Add(1)

		writer.Write(fmt.Sprintf("Trying to connect: %s - timeout = %v seconds", address, timeout))
		go func() {
			for true {

				if protocol=="tcp" {
					_, err := net.Dial("tcp", address)
					if err == nil {
						writer.Write(fmt.Sprintf("Connection OK : %s:%s", protocol, address))
						waitGroup.Add(-1)
						break
					} else if time.Now().Sub(start).Seconds() > timeout {
						writer.Write(fmt.Sprintf("%s : %s", err.Error(), address))
						waitGroup.Add(-1)
						failure = true
						break
					}
				} else {
					endpoint :=  protocol +"://" + address +"/"+ path
					resp, err := http.Get(endpoint)
					if err==nil && resp.StatusCode==successcode {
						writer.Write(fmt.Sprintf("Connection OK : %s", endpoint))
						waitGroup.Add(-1)
						break
					} else if time.Now().Sub(start).Seconds() > timeout {
						writer.Write(fmt.Sprintf("%s : %s", err.Error(), address))
						waitGroup.Add(-1)
						failure = true
						break
					}
				}

				time.Sleep(time.Millisecond * time.Duration(*pollinterval))
			}
		}()
	}

	waitGroup.Wait()

	if failure {
		os.Exit(*exitCodeOnConnectFail)
	}

	os.Exit(*exitCodeOnConnectOk)
}

func (writer outputWriter) Write(message string) {
	if !writer.QuietMode {
		fmt.Println(message)
	}
}

func showUsage() {
	fmt.Println("Switches (prefix with -)")
	fmt.Println("\t targets: comma separated protocol host, port, timeoutSeconds, path and successcode (last two for http/s only) e.g. localhost:8080:20,localhost:80:5")
	fmt.Println("\t globalTimeout: global timeout in seconds")
	fmt.Println("\t poll-interval: poll interval in milliseconds, defaults to 250")
	fmt.Println("\t exit-code-on-connect: defaults to 0")
	fmt.Println("\t exit-code-on-fail: defaults to 1")
	fmt.Println("\t quiet-mode: suppress all output, default true")
}
