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
 *     Initial: 2018/05/11        Chen Yanchen
 */

package issues

import (
	"context"

	"github.com/google/go-github/github"
	"golang.org/x/oauth2"
)

type serviceProvider struct{}

var Service *serviceProvider

// ListByAuthor get issues list.
func (*serviceProvider) List(token, owner, repo, creator string) ([]*github.Issue, error) {
	ctx := context.Background()
	tokenService := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: token},
	)
	tokenClient := oauth2.NewClient(ctx, tokenService)

	client := github.NewClient(tokenClient)

	opt := &github.IssueListByRepoOptions{
		Creator: creator,
		State:   "open",
	}

	list, _, err := client.Issues.ListByRepo(ctx, owner, repo, opt)
	if err != nil {
		return nil, err
	}

	return list, nil
}

// Get issues.
func (sp *serviceProvider) Get(token, owner, repo string, num int) (*github.Issue, error) {
	ctx := context.Background()
	tokenService := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: token},
	)
	tokenClient := oauth2.NewClient(ctx, tokenService)

	client := github.NewClient(tokenClient)

	issue, _, err := client.Issues.Get(ctx, owner, repo, num)
	if err != nil {
		return nil, err
	}

	return issue, nil
}
