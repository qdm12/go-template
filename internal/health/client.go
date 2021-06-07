package health

import (
	"context"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"
)

func IsClientMode(args []string) bool {
	return len(args) > 1 && args[1] == "healthcheck"
}

type Client interface {
	Query(ctx context.Context, address string) error
}

type client struct {
	*http.Client
}

func NewClient() Client {
	const timeout = 5 * time.Second
	return &client{
		Client: &http.Client{Timeout: timeout},
	}
}

var (
	ErrQuery     = errors.New("cannot query health server")
	ErrUnhealthy = errors.New("unhealthy")
)

// Query sends an HTTP request to the other instance of
// the program, and to its internal healthcheck server.
func (c *client) Query(ctx context.Context, address string) error {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, "http://"+address, nil)
	if err != nil {
		return fmt.Errorf("%w: %s", ErrQuery, err)
	}
	resp, err := c.Do(req)
	if err != nil {
		return fmt.Errorf("%w: %s", ErrQuery, err)
	} else if resp.StatusCode == http.StatusOK {
		return nil
	}

	b, err := ioutil.ReadAll(resp.Body)
	defer resp.Body.Close()
	if err != nil {
		return fmt.Errorf("%w: %s: %s", ErrUnhealthy, resp.Status, err)
	}
	return fmt.Errorf("%w: %s", ErrUnhealthy, string(b))
}
