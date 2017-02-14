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
	flag.Parse()

	if *address == "" {
		fmt.Println("Please supply address, e.g. -address=localhost:8080")
		os.Exit(1)
	}

	start := time.Now()

	for true {

		_, err := net.Dial("tcp", *address)

		if (err==nil) {
			fmt.Println("Connection OK")
			os.Exit(0)
		} else if (time.Now().Sub(start).Seconds()>*timeout) {
			fmt.Println(err.Error())
			os.Exit(1)
		}

		time.Sleep(time.Millisecond*time.Duration(*pollinterval))
	}
}
