/*
 * Revision History:
 *     Initial: 2018/07/04        Li Zebang
 */

package pool

import (
	"context"
	"errors"
	"strconv"
	"sync"
	"time"

	"github.com/google/go-github/github"
	"golang.org/x/oauth2"
)

const (
	defaultRemaining = -1

	// HeaderXRateLimitRemaining -
	HeaderXRateLimitRemaining = "X-RateLimit-Remaining"
	// HeaderXRateLimitReset -
	HeaderXRateLimitReset = "X-RateLimit-Reset"
)

// Token -
type Token struct {
	Tag   string
	Token string
}

// Client -
type Client struct {
	*github.Client

	isUsed    bool
	remaining int

	lock  *sync.Mutex
	reset time.Time
	tag   string
}

func client(token *Token) (*Client, error) {
	if token == nil {
		return nil, errors.New("token cannot be nil")
	}

	ctx := context.Background()

	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: token.Token},
	)

	tc := oauth2.NewClient(ctx, ts)

	return &Client{
		Client: github.NewClient(tc),

		lock:      &sync.Mutex{},
		remaining: defaultRemaining,
		tag:       token.Tag,
	}, nil
}

// HandleResponse -
func (c *Client) HandleResponse(resp *github.Response) error {
	c.lock.Lock()
	defer c.lock.Unlock()

	remaining, err := strconv.Atoi(resp.Header.Get(HeaderXRateLimitRemaining))
	if err != nil {
		return err
	}
	c.remaining = remaining

	if c.isUsed {
		return nil
	}

	reset, err := strconv.ParseInt(resp.Header.Get(HeaderXRateLimitReset), 10, 64)
	if err != nil {
		return err
	}
	c.reset = time.Unix(reset, 0)

	return nil
}

// Reset -
func (c *Client) Reset() {
	c.isUsed = false
	c.remaining = defaultRemaining
}
