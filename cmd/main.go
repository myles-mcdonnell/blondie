package main

import (
	"flag"
	"fmt"
	"github.com/myles-mcdonnell/blondie"
	"os"
	"strconv"
	"strings"
	"time"
)

func main() {

	targets := flag.String("targets", "", "comma separated protocol:address:port:timeoutSeconds:[path]:[successcode], e.g. [tcp|http|https]:localhost:8080:60:[path]:[success_error_code] timeout is optional in which case the global timeout is used, final two arguments for htyp only")
	pollinterval := flag.Int("poll-interval", 250, "poll interval in milliseconds")
	exitCodeOnConnectOk := flag.Int("exit-code-on-connect", 0, "Exit code when connection is made")
	exitCodeOnConnectFail := flag.Int("exit-code-on-fail", 1, "Exit code when connection is not made")
	quietMode := flag.Bool("quiet-mode", false, "suppress all output")
	flag.Parse()

	if len(*targets) < 1 {
		showUsage()
		return
	}

	options := blondie.DefaultOptions()

	options.ExitCodeOnOk = *exitCodeOnConnectOk
	options.ExitCodeOnFail = *exitCodeOnConnectFail
	options.PollInterval = time.Millisecond * time.Duration(*pollinterval)
	options.QuietMode = *quietMode

	targetsRaw := strings.Split(*targets, ",")
	depChecks := make([]blondie.DepCheck, len(targetsRaw))

	for index, target := range targetsRaw {
		protoAddrPortAndTimeout := strings.Split(target, ":")

		var protocol = strings.ToLower(protoAddrPortAndTimeout[0])

		if protocol != "tcp" && protocol != "http" {
			fmt.Printf("Unrecognised protocol %s", protocol)
			os.Exit(*exitCodeOnConnectFail)
		}

		host := protoAddrPortAndTimeout[1]

		var err error
		var port int
		if port, err = strconv.Atoi(protoAddrPortAndTimeout[2]); err != nil {
			fmt.Println("Unable to parse port %s", protoAddrPortAndTimeout[2])
			os.Exit(*exitCodeOnConnectFail)
		}

		to, err := strconv.Atoi(protoAddrPortAndTimeout[3])
		if err != nil {
			panic(err)
		}

		timeout := time.Duration(float64(to)) * time.Millisecond

		switch protocol {
		case "tcp":
			depChecks[index] = blondie.NewTcpCheck(host, port, timeout)
			break
		case "http":
			var path string
			var successcodes []int
			if len(protoAddrPortAndTimeout) > 4 {
				path = protoAddrPortAndTimeout[4]
			}

			if len(protoAddrPortAndTimeout) > 5 {
				successcode, err := strconv.Atoi(protoAddrPortAndTimeout[5])
				if err != nil {
					fmt.Println("Can not parse success code to int %s", protoAddrPortAndTimeout[5])
					os.Exit(*exitCodeOnConnectFail)
				} else {
					successcodes = []int{successcode}
				}
			} else {
				successcodes = []int{}
			}
			depChecks[index] = blondie.NewHttpCheck(host, port, timeout, path, successcodes)
			break
		default:
			fmt.Println("Unsupported protocol %s", protocol)
			os.Exit(*exitCodeOnConnectFail)
		}

	}

	if !blondie.WaitForDeps(depChecks, options) {
		os.Exit(*exitCodeOnConnectFail)
	}

	os.Exit(*exitCodeOnConnectOk)
}

func showUsage() {
	fmt.Println("Switches (prefix with -)")
	fmt.Println("\t targets: comma separated protocol host, port, timeoutMilliseconds, path and successcode (last two for http/s only) e.g. tcp:localhost:8080:2000,http:localhost:80:5:health:200")
	fmt.Println("\t poll-interval: poll interval in milliseconds, defaults to 250")
	fmt.Println("\t exit-code-on-connect: defaults to 0")
	fmt.Println("\t exit-code-on-fail: defaults to 1")
	fmt.Println("\t quiet-mode: suppress all output, default false")
}
