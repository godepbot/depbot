package sync

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"strings"
	"text/tabwriter"
	"time"

	"github.com/godepbot/depbot"
	"github.com/godepbot/depbot/internal/revision"
)

const (
	DepbotApiKey     = "DEPBOT_API_KEY"
	DepbotServerAddr = "DEPBOT_SERVER_ADDR"

	syncFlagApiKey       = "--api-key"
	syncFlagServerAddres = "--server-address"
)

var (
	// ErrorMissingApiKey is an error that will be returned if the API key is missing
	ErrorMissingApiKey error = fmt.Errorf("missing api key")
	// ErrorNoSyncDep is an error that will be returned if the dependencies could not be synchronized
	ErrorNoSyncDep error = fmt.Errorf("could not sync the dependencies")
)

type Command struct {
	finders []depbot.FinderFn

	stderr io.Writer
	stdout io.Writer
	stdin  io.Reader

	client *http.Client
}

func (c *Command) Name() string {
	return "sync"
}

func (c *Command) SetClient(client *http.Client) {
	c.client = client
}

func (c *Command) Main(ctx context.Context, pwd string, args []string) error {
	handleArgs(args)

	apiKey := os.Getenv(DepbotApiKey)
	if apiKey == "" {
		return ErrorMissingApiKey
	}

	hash, err := revision.FindLatestHash(pwd)
	if err != nil {
		return err
	}

	deps := depbot.Dependencies{}
	for _, df := range c.finders {
		dx, err := df(pwd)
		if err != nil {
			return err
		}

		deps = append(deps, dx...)
	}

	if len(deps) == 0 {
		return depbot.ErrorNoDependenciesFound
	}

	jm, err := json.Marshal(deps)
	if err != nil {
		return err
	}

	if c.client == nil {
		c.client = new(http.Client)
	}

	req, err := http.NewRequest(http.MethodPost, serverURL(), bytes.NewBuffer(jm))
	if err != nil {
		return fmt.Errorf("error creating new request %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %v", apiKey))
	req.Header.Set("X-Revision-Hash", strings.ReplaceAll(hash, "\n", ""))
	req.Header.Set("X-Timestamp", fmt.Sprintf("%v", time.Now().Unix()))

	resp, err := c.client.Do(req)
	if err != nil {
		return fmt.Errorf("error doing request %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		body, _ := ioutil.ReadAll(resp.Body)
		return fmt.Errorf("%v. error detail: %v", ErrorNoSyncDep.Error(), string(body))
	}

	defer resp.Body.Close()

	w := new(tabwriter.Writer)
	w.Init(c.stdout, 0, 0, 0, 0, 0)

	fmt.Fprintf(w, "%v dependencies synchronized.", len(deps))
	fmt.Fprintln(w)
	w.Flush()

	return nil
}

func (c *Command) SetIO(stderr io.Writer, stdout io.Writer, stdin io.Reader) {
	c.stderr = stderr
	c.stdout = stdout
	c.stdin = stdin
}

// NewCommand with the given finder function.
func NewCommand(finders ...depbot.FinderFn) *Command {
	return &Command{
		finders: finders,
	}
}

func serverURL() string {
	url := os.Getenv(DepbotServerAddr)
	if url == "" {
		url = "http://app.depbot.com/api/sync"
	}
	if !strings.Contains(url, "http") {
		url = fmt.Sprintf("http://%v", url)
	}
	return url
}

func handleArgs(args []string) {
	for _, arg := range args {
		flag := strings.Split(arg, "=")
		if len(flag) < 1 {
			continue
		}

		switch flag[0] {
		case syncFlagApiKey:
			os.Setenv(DepbotApiKey, flag[1])
		case syncFlagServerAddres:
			os.Setenv(DepbotServerAddr, flag[1])
		default:
			continue
		}
	}
}
