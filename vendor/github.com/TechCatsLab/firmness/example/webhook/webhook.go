/*
 * Revision History:
 *     Initial: 2018/05/25        Li Zebang
 */

package main

import (
	"fmt"

	"github.com/TechCatsLab/firmness/slack"
)

func main() {

	// You can add a Incoming WebHooks via
	// https://my.slack.com/services/new/incoming-webhook/
	webhook := slack.NewWebhook("webhook-link")

	err := webhook.PostMessage("Warning Message", "backend", "middleware")
	if err != nil {
		fmt.Println(err)
	}
}
