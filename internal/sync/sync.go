package sync

import (
	"bytes"
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"time"

	"github.com/godepbot/depbot"
	"github.com/godepbot/depbot/internal/revision"
)

const (
	DepbotApiKey     = "DEPBOT_API_KEY"
	DepbotServerAddr = "DEPBOT_SERVER_ADDR"

	// The default server address, could be changed with flags/env
	defaultServerAddress = "https://app.depbot.com/api/sync"
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

	client  *http.Client
	flagSet *flag.FlagSet

	apiKey        string
	serverAddress string

	revisionFinderHash   func(string) (string, error)
	revisionFinderBranch func(string) (string, error)
}

func (c *Command) Name() string {
	return "sync"
}

func (c *Command) HelpText() string {
	return "Synchronizes the dependencies with the server. The server can be configured with the \n\t\tDEPBOT_SERVER_ADDR environment variable. It requires a repo API key. See flags for more info."
}

func (c *Command) SetClient(client *http.Client) {
	c.client = client
}

func (c *Command) ParseFlags(args []string) (*flag.FlagSet, error) {
	_ = c.flagSet.Parse(args)

	// Applying env if apiKey is empty
	if c.apiKey == "" {
		c.apiKey = os.Getenv(DepbotApiKey)
	}

	// Applying env if passed value was empty.
	// IMPORTANT: We acknowledge that the server address flag will be
	// overridden by the env variable if it is set to be the default.
	if es := os.Getenv(DepbotServerAddr); c.serverAddress == defaultServerAddress && es != "" {
		c.serverAddress = es
	}

	return c.flagSet, nil
}

func (c *Command) Main(ctx context.Context, pwd string, args []string) error {
	if c.apiKey == "" {
		return ErrorMissingApiKey
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

	req, err := http.NewRequest(http.MethodPost, c.serverAddress, bytes.NewBuffer(jm))
	if err != nil {
		return fmt.Errorf("error creating new request %w", err)
	}

	// Finding current revision hash
	hash, err := c.revisionFinderHash(pwd)
	if err != nil {
		return err
	}

	branch, err := c.revisionFinderBranch(pwd)
	if err != nil {
		return err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %v", c.apiKey))
	req.Header.Set("X-Revision-Hash", hash)
	req.Header.Set("X-Revision-Branch", branch)
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

	fmt.Fprintf(c.stdout, "%v dependencies synchronized.", len(deps))

	return nil
}

func (c *Command) SetIO(stderr io.Writer, stdout io.Writer, stdin io.Reader) {
	c.stderr = stderr
	c.stdout = stdout
	c.stdin = stdin
}

// NewCommand with the given finder function.
func NewCommand(finders ...depbot.FinderFn) *Command {
	c := &Command{
		finders: finders,
		apiKey:  os.Getenv(DepbotApiKey),

		client:               http.DefaultClient,
		revisionFinderHash:   revision.FindLatestHash,
		revisionFinderBranch: revision.FindBranchName,
	}

	// Initializing the FlagSet that will pull the api-key and the server-address
	fls := flag.NewFlagSet(c.Name(), flag.ContinueOnError)
	fls.StringVar(&c.apiKey, "api-key", "", "[required] The API key for the repo. Can be specified with the DEPBOT_API_KEY environment variable.")
	fls.StringVar(&c.serverAddress, "server-address", defaultServerAddress, "The server address. Can be specified with the DEPBOT_SERVER_ADDR environment variable.")
	fls.SetOutput(bytes.NewBuffer([]byte{}))
	fls.Usage = func() {}

	c.flagSet = fls

	return c
}

// WithRevisionFinder is Useful for testing purposes so we can replace the
// Revision finder
func (s *Command) WithRevisionFinder(finderHash, finderBranch func(string) (string, error)) {

	s.revisionFinderHash = finderHash
	s.revisionFinderBranch = finderBranch
}
