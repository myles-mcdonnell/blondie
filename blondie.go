package main

import (
	"os"
	"fmt"
	"net"
	"time"
	"flag"
	"strings"
	"sync"
	"strconv"
)

type outputWriter struct {
	QuietMode bool
}

func main() {

	targets := flag.String("targets", "", "comma separated address:port:timeoutSeconds, e.g. localhost:8080:60 timeout is optional in whichcase the global timeout is used")
	globalTimeout := flag.Float64("globalTimeout", 5, "a global timeoput")
	pollinterval := flag.Int("poll-interval", 250, "poll interval in milliseconds")
	exitCodeOnConnectOk := flag.Int("exit-code-on-connect", 0, "Exit code when connection is made")
	exitCodeOnConnectFail := flag.Int("exit-code-on-fail", 1, "Exit code when connection is not made")
	quietMode := flag.Bool("quiet-mode", false, "suppress all output")
	flag.Parse()

	var waitGroup = &sync.WaitGroup{}
	var writer = &outputWriter{QuietMode:*quietMode}
	var failure bool

	for _, target := range strings.Split(*targets, ",") {
		addrPortAndTimeout := strings.Split(target, ":")

		var address = addrPortAndTimeout[0]+":"+addrPortAndTimeout[1]

		if address == "" {
			fmt.Println("Please supply targets, e.g. -targets=localhost:8080;60,google.com:80;30")
			os.Exit(*exitCodeOnConnectFail)
		}

		var timeout float64
		if len(addrPortAndTimeout)<3 {
			timeout = *globalTimeout
		} else {
			var err error
			to, err := strconv.Atoi(addrPortAndTimeout[2])
			if err!=nil {
				panic(err)
			}

			timeout = float64(to)
		}

		start := time.Now()

		waitGroup.Add(1)

		writer.Write(fmt.Sprintf("Trying to connect: %s - timeout = %v seconds", address, timeout))
		go func(){
			for true {
				_, err := net.Dial("tcp", address)

				if (err==nil) {
					writer.Write(fmt.Sprintf("Connection OK : %s", address))
					waitGroup.Add(-1)
					break//os.Exit(*exitCodeOnConnectOk)
				} else if (time.Now().Sub(start).Seconds()>timeout) {
					writer.Write(fmt.Sprintf("%s : %s", err.Error(), address))
					waitGroup.Add(-1)
					failure = true
					break//os.Exit(*exitCodeOnConnectFail)
				}

				time.Sleep(time.Millisecond*time.Duration(*pollinterval))
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
