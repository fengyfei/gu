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
 *     Initial: 2018/02/28        Li Zebang
 */

package crawler

import "fmt"

// Data -
type Data interface {
	String() string
	File() (string, string, string)
	IsFile() bool
}

// DefaultData -
type DefaultData struct {
	// Source -
	Source string
	// Date -
	Date string
	// Title -
	Title string
	// URL -
	URL string
	// Text -
	Text string
	// FileType -
	FileType string
}

func (data *DefaultData) String() string {
	if data.IsFile() {
		return fmt.Sprintf("Source: %s\nDate: %s\nTitle: %s\nURL: %s\n", data.Source, data.Date, data.Title, data.URL)
	}
	return fmt.Sprintf("Source: %s\nDate: %s\nTitle: %s\nURL: %s\n%s", data.Source, data.Date, data.Title, data.URL, data.Text)
}

// File return title, filetype and content if data is a file type.
func (data *DefaultData) File() (title, filetype, content string) {
	if data.IsFile() {
		return data.Title, data.FileType, data.Text
	}
	return "", "", ""
}

// IsFile return true if field FileType is not "".
func (data *DefaultData) IsFile() bool {
	return data.FileType != ""
}
