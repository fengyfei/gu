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
 *     Initial: 2018/01/29        Li Zebang
 */

package segment

import (
	"fmt"
	"strconv"
	"strings"
)

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
			case "del":
				return text[0] + mddel(todo[1])
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
			return text[0] + mda(todo[0][attrIndex+1:]+">"+todo[1])
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
	var (
		li    string
		count = 1
		text  = []string{"", s}
	)

	for {
		text = strings.SplitN(text[1], "</ul>", 2)
		if strings.Index(text[0], "<ul>") == -1 {
			count--
			li += mdli(text[0], count)
		}
		if count == 0 {
			return li + parseMD(text[1])
		}
		text = strings.SplitN(text[0]+"</ul>"+text[1], "<ul>", 2)
		li += mdli(text[0], count)
		count++
		text = strings.SplitN(text[1], "</ul>", 2)
		li += mdli(text[0], count)
		count--
	}
}

// *-+ text
func mdli(s string, n int) string {
	var li string
	for i := 1; i < n; i++ {
		li += "  "
	}

	s = strings.Replace(s, "<li>", "", -1)
	s = strings.Replace(s, "</li>", "", -1)
	s = strings.Replace(s, "<p>", li+"* ", -1)
	s = strings.Replace(s, "</p>", "\n", -1)

	return parseMD(s)
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
	return "<u>" + parseMD(text[0]) + "</u>" + parseMD(text[1])
}

// ~~text~~
func mddel(s string) string {
	text := strings.SplitN(s, "</del>", 2)
	if len(text) == 1 {
		return parseMD(" ~~" + text[0] + "~~ ")
	}
	return parseMD(" ~~"+text[0]+"~~ ") + parseMD(text[1])
}

// ![text](url)
func mdfigure(s string) string {
	text := strings.SplitN(s, "</figure>", 2)
	if len(text[0]) > 13 {
		img := strings.SplitN(text[0], "\"", 3)[1]
		if len(text) == 1 {
			return "![](" + img + ")"
		}
		return "![](" + img + ")\n" + parseMD(text[1])
	}
	return parseMD(text[1])
}

// # text
func mdh(s string) string {
	var header = "\n#"
	text := strings.SplitN(s, "</h", 2)
	h, _ := strconv.Atoi(text[0][:1])
	for len(header) < h {
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
	text := strings.SplitN(s, "\">", 2)
	if text[1][:4] != "<a h" {
		text := strings.SplitN(s, "</a>", 2)
		a := strings.SplitN(text[0], "\"", 3)
		if parseMD(a[2][1:]) == "" {
			return parseMD(text[1])
		} else if len(text) == 1 {
			return "[" + parseMD(a[2][1:]) + "]" + "(" + a[1] + ")"
		}
		return "[" + parseMD(a[2][1:]) + "]" + "(" + a[1] + ")" + parseMD(text[1])
	}

	text = strings.SplitN(s, "</a>", 2)
	counter := strings.Count(text[0], "<a href=\"")
	text = strings.SplitN(s, "<a ", counter+1)

	var a string
	for i := 0; i < counter; i++ {
		a += "</a>"
	}

	s = strings.Replace(text[counter], a, "", 1)

	return parseMD("<a " + s)
}
