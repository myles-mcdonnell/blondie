package blondie

import (
	//"errors"
	"errors"
	"net"
	"net/http"
	"testing"
	"time"
)

var (
	tcpCheckOk = &tcpCheck{
		dial: func(string, string, time.Duration) (net.Conn, error) { return nil, nil },
		netCheck: netCheck{
			host:    "localhost",
			port:    80,
			timeout: 1 * time.Second},
	}

	tcpCheckFail = &tcpCheck{
		dial: func(string, string, time.Duration) (net.Conn, error) { return nil, errors.New("test") },
		netCheck: netCheck{
			host:    "localhost",
			port:    80,
			timeout: 1 * time.Second},
	}

	httpCheckNoCodeNoPathOk = &httpCheck{
		get: func(string) (*http.Response, error) { return &http.Response{StatusCode: 200}, nil },
		netCheck: netCheck{
			host:    "localhost",
			port:    80,
			timeout: 1 * time.Second},
	}
)

func Test_TCP_Ok(t *testing.T) {
	if !WaitForDeps([]DepCheck{tcpCheckOk}, DefaultOptions()) {
		t.Fail()
	}
}

func Test_HTTP_Ok_NoCodeNoPath(t *testing.T) {
	if !WaitForDeps([]DepCheck{httpCheckNoCodeNoPathOk}, DefaultOptions()) {
		t.Fail()
	}
}

func Test_HTTP_Ok_200_NoPath(t *testing.T) {
	check := &httpCheck{
		get: func(string) (*http.Response, error) { return &http.Response{StatusCode: 200}, nil },
		netCheck: netCheck{
			host:    "localhost",
			port:    80,
			timeout: 1 * time.Second},
		successCodes: []int{200},
	}

	if !WaitForDeps([]DepCheck{check}, DefaultOptions()) {
		t.Fail()
	}
}

func Test_HTTP_Fail_200_NoPath(t *testing.T) {
	check := &httpCheck{
		get: func(string) (*http.Response, error) { return &http.Response{StatusCode: 500}, nil },
		netCheck: netCheck{
			host:    "localhost",
			port:    80,
			timeout: 1 * time.Second},
		successCodes: []int{200},
	}

	if WaitForDeps([]DepCheck{check}, DefaultOptions()) {
		t.Fail()
	}
}

func Test_HTTP_Fail_200_WithPath(t *testing.T) {
	check := &httpCheck{
		get: func(string) (*http.Response, error) { return &http.Response{StatusCode: 500}, nil },
		netCheck: netCheck{
			host:    "localhost",
			port:    80,
			timeout: 1 * time.Second},
		successCodes: []int{200},
		path:         "test",
	}

	if WaitForDeps([]DepCheck{check}, DefaultOptions()) {
		t.Fail()
	}
}

func Test_HTTP_Ok_200204_NoPath(t *testing.T) {
	check := &httpCheck{
		get: func(string) (*http.Response, error) { return &http.Response{StatusCode: 204}, nil },
		netCheck: netCheck{
			host:    "localhost",
			port:    80,
			timeout: 1 * time.Second},
		successCodes: []int{200, 204},
		path:         "test",
	}

	if !WaitForDeps([]DepCheck{check}, DefaultOptions()) {
		t.Fail()
	}
}

func Test_HTTP_Ok_200_WithPath(t *testing.T) {
	check := &httpCheck{
		get: func(string) (*http.Response, error) { return &http.Response{StatusCode: 200}, nil },
		netCheck: netCheck{
			host:    "localhost",
			port:    80,
			timeout: 1 * time.Second},
		successCodes: []int{200},
		path:         "test",
	}

	if !WaitForDeps([]DepCheck{check}, DefaultOptions()) {
		t.Fail()
	}
}

func Test_TCP_Fail(t *testing.T) {
	if WaitForDeps([]DepCheck{tcpCheckFail}, DefaultOptions()) {
		t.Fail()
	}
}

func Test_TCP_OK_HTTP_OK(t *testing.T) {
	if !WaitForDeps([]DepCheck{tcpCheckOk, httpCheckNoCodeNoPathOk}, DefaultOptions()) {
		t.Fail()
	}
}

func Test_TCP_Fail_HTTP_OK(t *testing.T) {
	if WaitForDeps([]DepCheck{tcpCheckFail, httpCheckNoCodeNoPathOk}, DefaultOptions()) {
		t.Fail()
	}
}
