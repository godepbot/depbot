package sync

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"strings"
	"text/tabwriter"
	"time"

	"github.com/godepbot/depbot"
	"github.com/godepbot/depbot/internal/depbotserver"
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
	deps := depbot.Dependencies{}
	for _, df := range c.finders {
		dx, err := df(pwd)
		if err != nil {
			return err
		}

		deps = append(deps, dx...)
	}

	if len(deps) == 0 {
		fmt.Fprintln(c.stdout, "No dependendies found to sync")
		return nil
	}

	hash, err := revision.FindLatestHash(pwd)
	if err != nil {
		return err
	}

	jm, err := json.Marshal(deps)
	if err != nil {
		return err
	}

	if c.client == nil {
		c.client = new(http.Client)
	}

	key := depbotserver.EnvValueForKey(depbotserver.DEPBOT_API_KEY)
	url := depbotserver.EnvValueForKey(depbotserver.DEPBOT_SERVER_ADDR)

	client := depbotserver.DepBotClient{
		Client: c.client,
		Input: depbotserver.DepBotInput{
			Time: time.Now().Unix(),
			KEY:  key,
			URL:  url,
			Body: bytes.NewBuffer(jm),
			Hash: strings.ReplaceAll(hash, "\n", ""),
		},
	}

	resp, err := client.Post()
	if err != nil {
		fmt.Fprintln(c.stdout, "Could not sync the dependencies. Error detail: ", err)
		return err
	}

	if resp.StatusCode != 200 {
		body, _ := ioutil.ReadAll(resp.Body)
		fmt.Fprintln(c.stdout, "Could not sync the dependencies. Error detail: ", string(body))
		return err
	}

	defer resp.Body.Close()

	w := new(tabwriter.Writer)
	w.Init(c.stdout, 0, 0, 0, 0, 0)

	fmt.Fprintf(w, "%v dependencies synchronized.", len(deps))
	fmt.Fprintln(w)
	w.Flush()

	return nil
}

func Send(jm []byte) {
	panic("unimplemented")
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
