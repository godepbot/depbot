package depbotserver

import (
	"bytes"
	"fmt"
	"net/http"
	"os"
)

const (
	DEPBOT_SERVER_ADDR = "DEPBOT_SERVER_ADDR"
	DEPBOT_API_KEY     = "DEPBOT_API_KEY"
)

var defaultKeys = map[string]string{
	DEPBOT_SERVER_ADDR: "http://app.depbot.com/api/sync",
	DEPBOT_API_KEY:     "key",
}

type DepBotInput struct {
	KEY  string
	URL  string
	Hash string
	Body *bytes.Buffer
	Time int64
}

type DepBotClient struct {
	Client *http.Client
	Input  DepBotInput
}

func (db *DepBotClient) Post() (*http.Response, error) {
	req, err := db.NewRequest()
	if err != nil {
		return nil, fmt.Errorf("error initiating new request %w", err)
	}

	resp, err := db.Client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("error doing request %w", err)
	}

	return resp, err
}

func (db DepBotClient) NewRequest() (*http.Request, error) {
	req, err := http.NewRequest(http.MethodPost, db.Input.URL, db.Input.Body)
	if err != nil {
		return nil, fmt.Errorf("error creating new request %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %v", db.Input.KEY))
	req.Header.Set("X-Revision-Hash", db.Input.Hash)
	req.Header.Set("X-Timestamp", fmt.Sprintf("%v", db.Input.Time))

	return req, nil
}

func EnvValueForKey(variable string) string {
	val := os.Getenv(variable)
	if val != "" {
		return val
	}
	return defaultKeys[variable]
}
