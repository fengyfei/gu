/*
 * Revision History:
 *     Initial: 2018/05/24        Li Zebang
 */

package slack

import (
	"fmt"

	"github.com/nlopes/slack"
)

// Messager contains required information
type Messager struct {
	client  *slack.Client
	channel string
	Config  *slack.PostMessageParameters
}

// NewMessager creates a new messager for channel messaging.
func NewMessager(token string, channel string, config *slack.PostMessageParameters) (*Messager, error) {
	if token == "" {
		return nil, fmt.Errorf("[firmness, slack, messager] token cann't be null")
	}

	if channel == "" {
		return nil, fmt.Errorf("[firmness, slack, messager] channel cann't be null")
	}

	var messager = &Messager{
		client:  slack.New(token),
		channel: channel,
	}

	if config == nil {
		cfg := slack.NewPostMessageParameters()
		messager.Config = &cfg
	} else {
		messager.Config = config
	}

	return messager, nil
}

// PostMessage send message to the specified channel.
func (m *Messager) PostMessage(message string, labels ...string) error {
	var text string

	if labels == nil {
		text = fmt.Sprintf("Message: %s", message)
	} else {
		text = fmt.Sprintf("Labels: %v\nMessage: %s", labels, message)
	}

	_, _, err := m.client.PostMessage(m.channel, text, *m.Config)
	return err
}
