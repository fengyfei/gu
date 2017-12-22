/*
 * MIT License
 *
 * Copyright (c) 2017 SmartestEE Co., Ltd..
 *
 * Permission is hereby granted, free of charge, to any person obtaining a copy
 * of this software and associated documentation files (the "Software"), to deal
 * in the Software without restriction, including without limitation the rights
 * to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
 * copies of the Software, and to permit persons to whom the Software is
 * furnished to do so, subject to the following conditions:
 *
 * The above copyright notice and this permission notice shall be included in all
 * copies or substantial portions of the Software.
 *
 * THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
 * IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
 * FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
 * AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
 * LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
 * OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
 * SOFTWARE.
 */

/*
 * Revision History:
 *     Initial: 2017/08/08        Jia Chenhui
 */

package client

import (
	"context"
	"errors"
	"net/http"
	"sync"
	"time"

	GitHub "github.com/google/go-github/github"
	"golang.org/x/oauth2"
)

type rateLimitCategory uint8

const (
	coreCategory rateLimitCategory = iota
	searchCategory
	categories
)

type GHClient struct {
	Client     *GitHub.Client
	rateLimits [categories]Rate
	Manager    *ClientManager
	timer      *time.Timer
	rateMu     sync.Mutex
}

type Rate struct {
	Limit     int
	Remaining int
	Reset     time.Time
}

// newClient create client based on token.
func newClient(token string) (client *GHClient, err error) {
	if token == "" {
		client = new(GHClient)
		tokenSource := new(oauth2.TokenSource)
		if !client.init(*tokenSource) {
			err = errors.New("failed to create client")
			return nil, err
		}

		return client, nil
	}

	client = new(GHClient)
	tokenSource := oauth2.StaticTokenSource(&oauth2.Token{AccessToken: token})
	if !client.init(tokenSource) {
		err = errors.New("failed to create client")
		return nil, err
	}

	return client, nil
}

// init initializes the client, returns true if available, or returns false.
func (c *GHClient) init(tokenSource oauth2.TokenSource) bool {
	httpClient := oauth2.NewClient(oauth2.NoContext, tokenSource)
	ghClient := GitHub.NewClient(httpClient)
	c.Client = ghClient

	if !c.isValidToken(httpClient) {
		return false
	}

	if c.isLimited() {
		return false
	}

	return true
}

// isValidToken check if token is valid.
func (c *GHClient) isValidToken(httpClient *http.Client) bool {
	resp, err := c.makeRequest(httpClient)
	if err != nil {
		return false
	}

	err = GitHub.CheckResponse(resp)
	if _, ok := err.(*GitHub.TwoFactorAuthError); ok {
		return false
	}

	return true
}

// makeRequest sends an HTTP GET request and returns an HTTP response, following
// policy (such as redirects, cookies, auth) as configured on the client.
func (c *GHClient) makeRequest(httpClient *http.Client) (*http.Response, error) {
	req, err := c.Client.NewRequest("GET", "", nil)
	if err != nil {
		return nil, err
	}

	resp, err := httpClient.Do(req)
	if err != nil {
		return nil, err
	}

	return resp, nil
}

// isLimited check if the client is available.
func (c *GHClient) isLimited() bool {
	rate, _, err := c.Client.RateLimits(context.Background())
	if err != nil {
		return true
	}

	response := new(struct {
		Resource *GitHub.RateLimits `json:"resource"`
	})
	response.Resource = rate

	if response.Resource != nil {
		c.rateMu.Lock()
		defer c.rateMu.Unlock()
		if response.Resource.Core != nil {
			c.rateLimits[coreCategory].Limit = response.Resource.Core.Limit
			c.rateLimits[coreCategory].Remaining = response.Resource.Core.Remaining
			c.rateLimits[coreCategory].Reset = response.Resource.Core.Reset.Time
			return false
		}
		if response.Resource.Search != nil {
			c.rateLimits[searchCategory].Remaining = response.Resource.Search.Remaining
			c.rateLimits[searchCategory].Limit = response.Resource.Search.Limit
			c.rateLimits[searchCategory].Reset = response.Resource.Search.Reset.Time
			return false
		}
	}

	return true
}

// initTimer initialize client timers.
func (c *GHClient) initTimer(resp *GitHub.Response) {
	if resp != nil {
		timer := time.NewTimer((*resp).Reset.Time.Sub(time.Now()) + time.Second*2)
		c.timer = timer

		return
	}
}

// newClients create a client list based on tokens.
func newClients(tokens []string) []*GHClient {
	var clients []*GHClient

	for _, t := range tokens {
		client, err := newClient(t)
		if err != nil {
			continue
		}

		clients = append(clients, client)
	}

	return clients
}

type ClientManager struct {
	reclaim  chan *GHClient
	Dispatch chan *GHClient
}

// run reclaim the client.
func (r *ClientManager) run(done chan bool) {
	for {
		select {
		case v := <-r.reclaim:
			r.Dispatch <- v
		case <-done:
			break
		}
	}
	close(r.Dispatch)
}

// NewClientManager create a new client manager based on tokens.
func NewClientManager(tokens []string) *ClientManager {
	var rb *ClientManager = &ClientManager{
		reclaim:  make(chan *GHClient),
		Dispatch: make(chan *GHClient, len(tokens)),
	}

	clients := newClients(tokens)
	done := make(chan bool)
	defer func() {
		done <- true
	}()

	go rb.run(done)
	go func() {
		for _, c := range clients {
			if !c.isLimited() {
				rb.reclaim <- c
			}
		}
	}()

	return rb
}

// GetClient get a client.
func (m *ClientManager) GetClient() *GHClient {
	for {
		select {
		case c := <-m.Dispatch:
			return c
		}
	}
}

// PutClient put the client back to the manager.
// resp: The response returned when calling the client.
func PutClient(client *GHClient, resp *GitHub.Response) {
	client.initTimer(resp)

	<-client.timer.C
	select {
	case client.Manager.reclaim <- client:
	}
}
