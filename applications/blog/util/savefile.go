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
 *     Initial: 2018/04/03        Chen Yanchen
 */

package util

import (
	"encoding/base64"
	"io/ioutil"
	"os"
	"strconv"
	"time"

	"github.com/fengyfei/gu/libs/logger"
	"regexp"
	"strings"
)

const (
	filePath  = "./file/"
	imagePath = "img/"
	mdPath    = "md/"
	mdsuffix  = "md"
	imgsuffix = "jpg"
)

func nameImage() string {
	return strconv.FormatInt(time.Now().Unix(), 10) + "." + imgsuffix
}

func nameMd(name string) string {
	name = strings.Replace(name, " ", "", 32)
	return name + strconv.FormatInt(time.Now().Unix(), 10) + "." + mdsuffix
}

func checkDir(path string) error {
	_, err := os.Stat(path)

	if err != nil {
		if os.IsExist(err) {
			return err
		}
		err = os.MkdirAll(path, 0777)
		if err != nil {
			return err
		}
	}

	return nil
}

func SaveImage(base64Str string, path string) (string, error) {
	fileName := nameImage()
	path = filePath + imagePath + path

	img, err := base64.StdEncoding.DecodeString(base64Str)
	if err != nil {
		logger.Error(err)
		return "", err
	}

	err = checkDir(path)
	if err != nil {
		logger.Error("Check dir false:", err)
		return "", err
	}

	err = ioutil.WriteFile(path+fileName, img, 0777)
	if err != nil {
		logger.Error("WriteFile false:", err)
	}

	return path + fileName, err
}

func SaveMarkdown(content, name, path string) (string, error) {
	fileName := nameMd(name)
	path = filePath + mdPath + path

	err := checkDir(path)
	if err != nil {
		logger.Error("Check dir false:", err)
		return "", err
	}

	err = ioutil.WriteFile(path+fileName, []byte(content), 0777)
	if err != nil {
		logger.Error("WriteFile false:", err)
	}
	return path + fileName, err
}

func Delete(name string) bool {
	_, err := os.Stat(name)
	if err == nil || os.IsExist(err) {
		err = os.Remove(name)
		if err != nil {
			return false
		}
	}
	return true
}

func GetBrief(content string) string {
	reg, _ := regexp.Compile("\\<[\\S\\s]+?\\>")
	brief := reg.ReplaceAllString(content, "")
	return brief
}
