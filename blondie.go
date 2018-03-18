package blondie

import (
	"fmt"
	"net"
	"net/http"
	"sync"
	"time"
)

type Options struct {
	PollInterval   time.Duration
	ExitCodeOnOk   int
	ExitCodeOnFail int
	QuietMode      bool
	OutputWriter   func(string)
}

func DefaultOptions() *Options {
	return &Options{
		PollInterval:   time.Millisecond * 250,
		ExitCodeOnFail: 1,
		ExitCodeOnOk:   0,
		QuietMode:      false,
		OutputWriter:   func(msg string) { fmt.Println(msg) },
	}
}

type netCheck struct {
	host    string
	port    int
	timeout time.Duration
}

type tcpCheck struct {
	netCheck
	dial func(string, string) (net.Conn, error)
}

func NewTcpCheck(host string, port int, timeout time.Duration) DepCheck {
	return &tcpCheck{
		netCheck: netCheck{
			host:    host,
			port:    port,
			timeout: timeout,
		},
		dial: net.Dial,
	}
}

type httpCheck struct {
	netCheck
	path         string
	successCodes []int
	get          func(string) (*http.Response, error)
}

// NewHttpCheck creates a new DepCheck for a HTTP endpoint. Path may be an empty string and successCodes may be an empty slice in which case any reponse code will be considered a successful response
func NewHttpCheck(host string, port int, timeout time.Duration, path string, successCodes []int) DepCheck {
	return &httpCheck{
		netCheck: netCheck{
			host:    host,
			port:    port,
			timeout: timeout,
		},
		successCodes: successCodes,
		path:         path,
		get:          http.Get,
	}
}

type DepCheck interface {
	Try() bool
	Timeout() time.Duration
	Address() string
}

func (check *netCheck) Timeout() time.Duration {
	return check.timeout
}

func (check *httpCheck) Try() bool {
	endpoint := fmt.Sprintf("http://%s:%v/%s", check.host, check.port, check.path)
	resp, err := check.get(endpoint)

	if err == nil {
		if len(check.successCodes) == 0 {
			return true
		}

		for _, successCode := range check.successCodes {
			if successCode == resp.StatusCode {
				return true
			}
		}
	}

	return false
}

func (check *tcpCheck) Try() bool {
	address := fmt.Sprintf("%s:%v", check.host, check.port)
	_, err := check.dial("tcp", address)
	return err == nil
}

func (check *httpCheck) Address() string {
	return fmt.Sprintf("http://%s:%v/%s", check.host, check.port, check.path)
}

func (check *tcpCheck) Address() string {
	return fmt.Sprintf("tcp://%s:%v", check.host, check.port)
}

func (options Options) Write(message string) {
	if !options.QuietMode {
		options.OutputWriter(message)
	}
}

func WaitForDeps(deps []DepCheck, opts *Options) bool {

	var waitGroup = &sync.WaitGroup{}
	success := true

	waitGroup.Add(len(deps))
	for _, target := range deps {
		opts.Write(fmt.Sprintf("Trying to connect: %s - timeout = %v seconds", target.Address(), target.Timeout()))
		go func(target DepCheck) {
			start := time.Now()
			for true {
				if target.Try() {
					waitGroup.Done()
					break
				} else if time.Now().Sub(start) > target.Timeout() {
					opts.Write(fmt.Sprintf("Timeout : %s", target.Address()))
					success = false
					waitGroup.Done()
					break
				}

				time.Sleep(opts.PollInterval)
			}
		}(target)
	}

	waitGroup.Wait()

	return success
}

/*func main() {
	if protocol == "tcp" {
					_, err := net.Dial("tcp", address)
					if err == nil {
						writer.Write(fmt.Sprintf("Connection OK : %s:%s", protocol, address))
						waitGroup.Add(-1)
						break
					} else if time.Now().Sub(start).Seconds() > timeout {
						writer.Write(fmt.Sprintf("%s : %s", err.Error(), address))

						failure = true
						break
					}
				} else {
					endpoint := protocol + "://" + address + "/" + path
					resp, err := http.Get(endpoint)
					if err == nil && resp.StatusCode == successcode {
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

		if protocol != "tcp" && protocol != "http" && protocol != "https" {
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
		if protocol == "http" || protocol == "https" {
			path = protoAddrPortAndTimeout[4]
			var err error
			successcode, err = strconv.Atoi(protoAddrPortAndTimeout[5])
			if err != nil {
				fmt.Println("Can not parse success code to int %s", protoAddrPortAndTimeout[5])
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

				if protocol == "tcp" {
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
					endpoint := protocol + "://" + address + "/" + path
					resp, err := https.Get(endpoint)
					if err == nil && resp.StatusCode == successcode {
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
}*/
