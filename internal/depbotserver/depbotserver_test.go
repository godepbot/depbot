package depbotserver_test

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/godepbot/depbot/internal/depbotserver"
)

const (
	depSynchronized    = "dependencies synchronized"
	depNotSynchronized = "Could not sync the dependencies. Error detail"
	body               = `
	[
		{
		   "name": "github.com/PuerkitoBio/goquery",
		   "version": "v1.5.0",
		   "file": "go.mod",
		   "direct": true,
		   "kind": "library",
		   "license": "GNU V2"
		}
	]
	`
)

func mockEndPoint(w http.ResponseWriter, r *http.Request) {
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusUnprocessableEntity)
		w.Write([]byte(depNotSynchronized))
	}

	if strings.Contains(string(body), "version") {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(depSynchronized))
	} else {
		w.WriteHeader(http.StatusUnprocessableEntity)
		w.Write([]byte(depNotSynchronized))
	}
}

func Test_DepbotServer_Post_statusOK(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(mockEndPoint))
	defer server.Close()

	depbotClient := depbotserver.DepBotClient{
		Client: server.Client(),
		Input: depbotserver.DepBotInput{
			KEY:  "key xample",
			URL:  server.URL,
			Hash: "[hash xample]",
			Body: bytes.NewBuffer([]byte(body)),
			Time: time.Now().Unix(),
		},
	}

	res, err := depbotClient.Post()
	if err != nil {
		t.Fatalf("Got %v, but was expected %v", err, nil)
	}

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		t.Fatalf("Got %v, but was expected %v", err, nil)
	}

	if string(body) != depSynchronized {
		t.Fatalf("Got %v, but was expected %v", string(body), depSynchronized)
	}
}

func Test_DepbotServer_Post_statusUnprocesable(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(mockEndPoint))
	defer server.Close()

	depbotClient := depbotserver.DepBotClient{
		Client: server.Client(),
		Input: depbotserver.DepBotInput{
			KEY:  "key xample",
			URL:  server.URL,
			Hash: "[hash xample]",
			Body: bytes.NewBuffer([]byte("")),
			Time: time.Now().Unix(),
		},
	}

	res, err := depbotClient.Post()
	if err != nil {
		t.Fatalf("Got %v, but was expected %v", err, nil)
	}

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		t.Fatalf("Got %v, but was expected %v", err, nil)
	}

	if string(body) != depNotSynchronized {
		t.Fatalf("Got %v, but was expected %v", string(body), depNotSynchronized)
	}
}

func Test_Depbot_NewRequest(t *testing.T) {

	body := bytes.NewBuffer([]byte(body))

	input := depbotserver.DepBotInput{
		KEY:  "key xample",
		URL:  "http://127.0.0.1/api/sync",
		Hash: "[hash xample]",
		Body: body,
		Time: time.Now().Unix(),
	}

	depbotInput := depbotserver.DepBotClient{
		Input: input,
	}

	req, err := depbotInput.NewRequest()
	if err != nil {
		t.Fatalf("Got %v, but was expected %v", err.Error(), nil)
	}
	if req.Header.Get("Content-Type") != "application/json" {
		t.Fatalf("Got %v, but was expected %v", req.Header.Get("Content-Type"), "application/json")
	}
	if req.Header.Get("Authorization") != "Bearer "+input.KEY {
		t.Fatalf("Got %v, but was expected %v", req.Header.Get("Authorization"), "Bearer "+input.KEY)
	}
	if req.Header.Get("X-Revision-Hash") != input.Hash {
		t.Fatalf("Got %v, but was expected %v", req.Header.Get("X-Revision-Hash"), input.Hash)
	}
	if req.Header.Get("X-Timestamp") != fmt.Sprintf("%v", input.Time) {
		t.Fatalf("Got %v, but was expected %v", req.Header.Get("X-Timestamp"), input.Time)
	}
}
