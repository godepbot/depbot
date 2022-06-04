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
	"text/tabwriter"
	"time"

	"github.com/godepbot/depbot"
	"github.com/godepbot/depbot/internal/revision"
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
	apiKey := os.Getenv(depbot.EnvVariable_ApiKey)
	if apiKey == "" {
		return depbot.ErrorMissingApiKey
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

	url := os.Getenv(depbot.EnvVariable_ServerADDR)
	if url == "" {
		url = "http://app.depbot.com/api/sync"
	}

	if c.client == nil {
		c.client = new(http.Client)
	}

	req, err := http.NewRequest(http.MethodPost, url, bytes.NewBuffer(jm))
	if err != nil {
		return fmt.Errorf("error creating new request %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %v", apiKey))
	req.Header.Set("X-Revision-Hash", hash)
	req.Header.Set("X-Timestamp", fmt.Sprintf("%v", time.Now().Unix()))

	resp, err := c.client.Do(req)
	if err != nil {
		return fmt.Errorf("error doing request %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		body, _ := ioutil.ReadAll(resp.Body)
		return fmt.Errorf("%v. error detail: %v", depbot.ErrorNoSyncDep.Error(), string(body))
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
