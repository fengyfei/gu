/*
 * MIT License
 *
 * Copyright (c) 2018 SmartestEE Co., Ltd..
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
 *     Initial: 2018/02/10        Li Zebang
 */

package slack

import (
	"github.com/nlopes/slack"
)

type SlackClient struct {
	client *slack.Client
}

func NewClient(token string) *SlackClient {
	return &SlackClient{
		client: slack.New(token),
	}
}

func (sc *SlackClient) PostMessage(channel, text string) error {
	params := slack.NewPostMessageParameters()
	_, _, err := sc.client.PostMessage(channel, text, params)
	return err
}

func (sc *SlackClient) UploadFile(channel, title, filetype, content string) error {
	params := slack.FileUploadParameters{
		Title:    title,
		Filetype: filetype,
		Content:  content,
		Channels: []string{channel},
	}
	_, err := sc.client.UploadFile(params)
	return err
}
