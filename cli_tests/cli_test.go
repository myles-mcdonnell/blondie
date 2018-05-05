// +build cli_tests

package cli_tests

import (
	"log"
	"net/http"
	"os"
	"os/exec"
	"strconv"
	"testing"
)

func handler(w http.ResponseWriter, r *http.Request) {
	code, _ := strconv.Atoi(r.URL.Query().Get("code"))
	w.WriteHeader(code)
}

func TestMain(m *testing.M) {

	srv := &http.Server{Addr: ":8080"}
	http.HandleFunc("/", handler) // set router
	go func() {
		if err := srv.ListenAndServe(); err != nil {
			log.Fatal("ListenAndServe: ", err)
		}
	}()
	defer srv.Shutdown(nil)

	if err := exec.Command("go", "build", "-o=../artefacts/blondie", "../cmd/blondie/main.go").Run(); err != nil {
		panic(err)
	}

	os.Exit(m.Run())
}

func Test_HTTP_OK_NoPath_SingleCode(t *testing.T) {

	if err := exec.Command("../artefacts/blondie", "--targets=http:localhost:8080:3000:?code=200:200").Run(); err != nil {
		t.Log(err)
		t.Fail()
	}
}

func Test_HTTP_OK_NoPath_MultiCode(t *testing.T) {

	if err := exec.Command("../artefacts/blondie", "--targets=http:localhost:8080:3000:?code=200:200_204").Run(); err != nil {
		t.Log(err)
		t.Fail()
	}
}

func Test_HTTP_Fail_NoPath_SingleCode(t *testing.T) {

	if err := exec.Command("../artefacts/blondie", "--targets=http:localhost:8080:3000:?code=200:204").Run(); err == nil {
		t.Log(err)
		t.Fail()
	}
}

func Test_HTTP_Fail_NoPath_MultiCode(t *testing.T) {

	if err := exec.Command("../artefacts/blondie", "--targets=http:localhost:8080:3000:?code=200:203_204").Run(); err == nil {
		t.Log(err)
		t.Fail()
	}
}

func Test_HTTP_OK_TCP_OK(t *testing.T) {

	if err := exec.Command("../artefacts/blondie", "--targets=http:localhost:8080:3000:?code=200:200,tcp:localhost:8080:3000").Run(); err != nil {
		t.Log(err)
		t.Fail()
	}
}

func Test_HTTP_BadAddress_TCP_OK(t *testing.T) {

	if err := exec.Command("../artefacts/blondie", "--targets=http:nosuchaddress:8080:3000:?code=200:200,tcp:localhost:8080:3000").Run(); err == nil {
		t.Log(err)
		t.Fail()
	}
}
