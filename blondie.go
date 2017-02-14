package main

import (
	"os"
	"fmt"
	"net"
	"time"
	"flag"
)

func main() {

	timeout := flag.Float64("timeout", 5, "timeout in seconds")
	address := flag.String("address", "", "address and port, e.g. localhost:8080")
	pollinterval := flag.Int("poll-interval", 250, "poll interval in milliseconds")
	exitCodeOnConnectOk := flag.Int("exit-code-on-connect", 0, "Exit code when connection is made")
	exitCodeOnConnectFail := flag.Int("exit-code-on-fail", 1, "Exit code when connection is not made")
	quietMode := flag.Bool("quiet-mode", true, "suppress all output")
	flag.Parse()

	if *address == "" {
		fmt.Println("Please supply address, e.g. -address=localhost:8080")
		os.Exit(*exitCodeOnConnectFail)
	}

	start := time.Now()

	for true {

		_, err := net.Dial("tcp", *address)

		if (err==nil) {
			if !*quietMode {
				fmt.Println("Connection OK")
			}
			os.Exit(*exitCodeOnConnectOk)
		} else if (time.Now().Sub(start).Seconds()>*timeout) {
			if !*quietMode {
				fmt.Println(err.Error())
			}
			os.Exit(*exitCodeOnConnectFail)
		}

		time.Sleep(time.Millisecond*time.Duration(*pollinterval))
	}
}
