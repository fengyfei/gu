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
 *     Initial: 2018/03/29        Tong Yuehong
 */

package util

import (
	"io/ioutil"
	"os"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/fengyfei/gu/libs/logger"
)

var (
	content = "./content/"
	image   = "./image/"
	suffix  = "txt"
	Content = 0
	Image   = 1
)

func fileName(userID uint32, diff int) string {
	timestamp := time.Now().Unix()

	time := time.Unix(timestamp, 0).Format("2006-01-02 03:04:05 PM")
	time = strings.Replace(time, " ", "", 2)

	id := strconv.FormatUint(uint64(userID), 10)

	if diff == Content {
		return content + time + id + "." + suffix
	} else {
		return image + time + id + "." + suffix
	}
}

func Save(userID uint32, content string, diff int) (string, error) {
	fileName := fileName(userID, diff)

	err := ioutil.WriteFile(fileName, []byte(content), 0777)
	if err != nil {
		logger.Error(err)
	}

	return fileName, err
}

func DeleteFile(userID uint32, diff int) bool {
	path := fileName(userID, diff)
	_, err := os.Stat(path)

	if err == nil || os.IsExist(err) {
		err = os.Remove(path)
		if err != nil {
			return false
		}
		return true
	}

	return os.IsNotExist(err)
}

func GetBrief(content string) string {
	re, _ := regexp.Compile("\\<[\\S\\s]+?\\>")
	brief := re.ReplaceAllString(content, "")

	return brief
}
