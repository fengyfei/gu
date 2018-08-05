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

	// Both App Bots and Custom Bots can be used here.
	// More Information: https://api.slack.com/bot-users
	//
	// App Bots: You can create a Slack app and add a bot in the
	// management dashboard.
	// Add a Slack app: https://api.slack.com/apps/new
	//
	// Custom Bots: click https://my.slack.com/services/new/bot and
	// you can create a custom bot.
	messager, err := slack.NewMessager("bot-token", "channel", nil)
	if err != nil {
		fmt.Println(err)
	}

	err = messager.PostMessage("Warning Message", "backend", "middleware")
	if err != nil {
		fmt.Println(err)
	}
}
