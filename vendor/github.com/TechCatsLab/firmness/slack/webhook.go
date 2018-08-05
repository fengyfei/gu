/*
 * Revision History:
 *     Initial: 2018/05/25        Li Zebang
 */

package slack

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"
)

var httpClient = &http.Client{}

// Webhook contains required information
type Webhook struct {
	url    string
	client *http.Client
}

// NewWebhook return a incoming webhook for sending data
// into Slack in real-time.
func NewWebhook(url string) *Webhook {
	return &Webhook{
		url:    url,
		client: httpClient,
	}
}

// PostMessage send message to the specified channel.
// This channel is determined when you create a webhook.
func (w *Webhook) PostMessage(message string, labels ...string) error {
	var text string

	if labels == nil {
		text = fmt.Sprintf("Message: %s", message)
	} else {
		text = fmt.Sprintf("Labels: %v\nMessage: %s", labels, message)
	}

	reqBody, err := json.Marshal(basicData{text})
	if err != nil {
		return err
	}

	req, err := http.NewRequest("POST", w.url, strings.NewReader(string(reqBody)))
	if err != nil {
		return err
	}
	defer req.Body.Close()

	resp, err := w.client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusTooManyRequests {
		retry, err := strconv.ParseInt(resp.Header.Get("Retry-After"), 10, 64)
		if err != nil {
			return err
		}
		return &RateLimitedError{time.Duration(retry) * time.Second}
	}

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("[firmness, slack, webhook] slack server error: %s", resp.Status)
	}

	return nil
}

type basicData struct {
	Text string `json:"text"`
}

// RateLimitedError implements error interface
type RateLimitedError struct {
	RetryAfter time.Duration
}

func (e *RateLimitedError) Error() string {
	return fmt.Sprintf("Slack rate limit exceeded, retry after %s", e.RetryAfter)
}
