/*
 * Revision History:
 *     Initial: 2018/05/26        Li Zebang
 */

package mail

import (
	"fmt"
	"net/mail"
	"net/smtp"

	"github.com/TechCatsLab/firmness/checker"
)

// Client contains required sender information
type Client struct {
	config *Config
}

// NewClient create a sender client.
func NewClient(config *Config) (*Client, error) {
	if config == nil {
		return nil, fmt.Errorf("config cann't be nil")
	}

	if !checker.IsEmail(config.From) {
		return nil, fmt.Errorf("the account's email %s is invalid", config.From)
	}

	if !checker.IsEmail(config.Credentials.Username) || config.Credentials.Password == "" {
		return nil, fmt.Errorf("the account's email %s is invalid", config.From)
	}

	return &Client{config}, nil
}

// PostMessage send the message to the specified account.
func (c *Client) PostMessage(subject, message string, labels []string) error {
	var auth = smtp.PlainAuth("", c.config.Credentials.Username, c.config.Credentials.Password, c.config.Host)

	addrlist, err := mail.ParseAddressList(c.config.To)
	if err != nil {
		return err
	}

	to := make([]string, 0)
	for _, addr := range addrlist {
		to = append(to, addr.Address)
	}

	msg := fmt.Sprintf("From: %s\r\nTo: %s\r\nSubject: %s\r\n\r\nLabels: %v\nMessage: %s", c.config.From, c.config.To, subject, labels, message)

	return smtp.SendMail(
		fmt.Sprintf("%s:%s", c.config.Host, c.config.Port),
		auth,
		c.config.From,
		to,
		[]byte(msg),
	)
}
