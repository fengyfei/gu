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
 *     Initial: 2018/03/27        Shi Ruitao
 */

package user

import (
	"encoding/base64"
	"errors"
	"io/ioutil"
	"os"
	"strconv"
	"strings"

	"github.com/fengyfei/gu/libs/logger"
)

var picturePath = "./img/"

func checkDir(path string) error {
	_, err := os.Stat(path)

	if err != nil {
		if os.IsNotExist(err) {
			err = os.MkdirAll(path, 0777)

			if err != nil {
				return err
			}
		}
	}

	return err
}

func SavePicture(base64Str string, pathPrefix string, id uint32) (string, error) {
	if !strings.Contains(base64Str, "base64") || !strings.Contains(base64Str, "image") {
		return "", errors.New("unvalid image base64 string")
	}

	slice := strings.Split(base64Str, ",")
	suffix := string([]byte(slice[0])[11 : len(slice[0])-7]) // picture format suffix

	byteData, err := base64.StdEncoding.DecodeString(slice[1])
	if err != nil {
		logger.Error(err)

		return "", err
	}

	fileName := strconv.FormatInt(int64(id), 10) + "." + suffix
	path := picturePath + pathPrefix
	err = checkDir(path)
	if err != nil {
		logger.Error(err)

		return "", err
	}

	err = ioutil.WriteFile(path+fileName, byteData, 0777)
	if err != nil {
		logger.Error(err)
	}

	return path + fileName, err
}
