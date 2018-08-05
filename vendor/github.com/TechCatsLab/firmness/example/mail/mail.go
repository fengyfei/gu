/*
 * Revision History:
 *     Initial: 2018/05/26        Li Zebang
 */

package main

import (
	"fmt"

	"github.com/TechCatsLab/firmness/mail"
)

func main() {

	// Please make sure to enable smtp service before use.
	config := &mail.Config{
		From: "xxx@163.com",
		To:   "xxx@gmail.com, xxx@outlook.com",
		Host: "smtp.163.com",
		Port: "25",
		Credentials: mail.Credentials{
			Username: "xxx@163.com",
			Password: "xxx",
		},
	}

	client, err := mail.NewClient(config)
	if err != nil {
		fmt.Println(err)
	}

	err = client.PostMessage("subject", "message", []string{"backend", "middleware"})
	if err != nil {
		fmt.Println(err)
	}
}
