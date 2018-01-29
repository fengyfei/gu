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
 *     Initial: 2017/01/29        Li Zebang
 */

package segment

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/fengyfei/gu/libs/logger"
)

// &quot;

func parseMD(s string) string {
	for {
		if strings.Index(s, "<") == -1 {
			return s
		}

		text := strings.SplitN(s, "<", 2)
		todo := strings.SplitN(text[1], ">", 2)
		attrIndex := strings.Index(todo[0], " ")

		if attrIndex == -1 {
			switch todo[0] {
			case "p":
				return text[0] + mdp(todo[1])
			case "code":
				return text[0] + mdcode(todo[1])
			case "em":
				return text[0] + mdem(todo[1])
			case "strong":
				return text[0] + mdstrong(todo[1])
			case "blockquote":
				return text[0] + mdblockquote(todo[1])
			case "ul":
				return text[0] + mdul(todo[1])
			case "ol":
				return text[0] + mdol(todo[1])
			case "u":
				return text[0] + mdu(todo[1])
			case "figure":
				return text[0] + mdfigure(todo[0][attrIndex+1:]+todo[1])
			}
		}

		if todo[0][:1] == "h" {
			return text[0] + mdh(todo[0][1:2]+todo[1])
		}

		switch todo[0][:attrIndex] {
		case "pre":
			return text[0] + mdpre(todo[0][attrIndex+1:]+todo[1])
		case "a":
			return text[0] + mda(todo[0][attrIndex+1:]+todo[1])

		}
	}
}

// \n
func mdp(s string) string {
	text := strings.SplitN(s, "</p>", 2)
	if len(text) == 1 {
		return "\n" + parseMD(text[0]) + "\n"
	}
	return "\n" + parseMD(text[0]) + "\n" + parseMD(text[1])
}

// `text`
func mdcode(s string) string {
	text := strings.SplitN(s, "</code>", 2)
	if len(text) == 1 {
		return parseMD(" `" + text[0] + "` ")
	}
	return parseMD(" `"+text[0]+"` ") + parseMD(text[1])
}

// *text*
func mdem(s string) string {
	text := strings.SplitN(s, "</em>", 2)
	if len(text) == 1 {
		return parseMD(" *" + text[0] + "* ")
	}
	return parseMD(" *"+text[0]+"* ") + parseMD(text[1])
}

// **text**
func mdstrong(s string) string {
	text := strings.SplitN(s, "</strong>", 2)
	if len(text) == 1 {
		return parseMD(" **" + text[0] + "** ")
	}
	return parseMD(" **"+text[0]+"** ") + parseMD(text[1])
}

// > text
func mdblockquote(s string) string {
	text := strings.SplitN(s, "</blockquote>", 2)
	blockquote := strings.Replace(text[0], "<p>", "> ", -1)
	blockquote = strings.Replace(blockquote, "</p>", "\n", -1)
	if len(text) == 1 {
		return parseMD(blockquote)
	}
	return parseMD(blockquote) + parseMD(text[1])
}

// *-+ text
func mdul(s string) string {
	text := strings.SplitN(s, "</ul>", 2)
	ul := strings.Replace(text[0], "<li>", "* ", -1)
	ul = strings.Replace(ul, "</li>", "\n", -1)
	for _, v := range []string{"<p>", "</p>"} {
		ul = strings.Replace(ul, v, "", -1)
	}
	if len(text) == 1 {
		return parseMD(ul)
	}
	return parseMD(ul) + parseMD(text[1])
}

// 1. text
func mdol(s string) string {
	var str string
	text := strings.SplitN(s, "</ol>", 2)
	ol := strings.Split(strings.Replace(text[0], "</p></li>", "", -1), "<li><p>")
	for k, v := range ol[1:] {
		str += fmt.Sprintf("%d. %s\n", k+1, v)
	}
	if len(text) == 1 {
		return parseMD(str)
	}
	return parseMD(str) + parseMD(text[1])
}

// <u>text</u>
func mdu(s string) string {
	text := strings.SplitN(s, "</u>", 2)
	if len(text) == 1 {
		return "<u>" + text[0] + "</u>"
	}
	return "<u>" + parseMD(text[0]) + "</u>\n" + parseMD(text[1])
}

// ![text](url)
func mdfigure(s string) string {
	text := strings.SplitN(s, "</figure>", 2)
	img := strings.SplitN(text[0], "\"", 3)[1]
	if len(text) == 1 {
		return "![](" + img + ")"
	}
	return "![](" + img + ")\n" + parseMD(text[1])
}

// # text
func mdh(s string) string {
	var header string
	text := strings.SplitN(s, "</h", 2)
	h, err := strconv.Atoi(text[0][:1])
	if err != nil {
		logger.Error("error in parsing <h>", err)
		errorPipe <- err
	}
	for len(header) <= h {
		header += "#"
	}
	if len(text) == 1 {
		return parseMD(header + " " + text[0][1:])
	}
	return parseMD(header+" "+text[0][1:]) + "\n" + parseMD(strings.SplitN(text[1], ">", 2)[1])
}

// ``` language```
func mdpre(s string) string {
	text := strings.SplitN(s, "</pre>", 2)
	language := strings.SplitN(text[0], "\"", 3)
	start := strings.Index(language[2], "<code>")
	end := strings.Index(language[2], "</code>")
	if len(text) == 1 {
		return "\n``` " + language[1] + "\n" + language[2][start+6:end] + "\n```\n"
	}
	return "\n``` " + language[1] + "\n" + language[2][start+6:end] + "\n```\n" + parseMD(text[1])
}

// [text](url)
func mda(s string) string {
	text := strings.SplitN(s, "</a>", 2)
	a := strings.SplitN(text[0], "\"", 3)
	if len(text) == 1 {
		return "[" + parseMD(a[2]) + "]" + "(" + a[1] + ")"
	}
	return "[" + parseMD(a[2]) + "]" + "(" + a[1] + ")" + parseMD(text[1])
}
