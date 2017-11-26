/*
 * MIT License
 *
 * Copyright (c) 2017 SmartestEE Co., Ltd..
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
 *     Initial: 2017/11/26        Wang RiYu
 */

package util

import (
  "encoding/base64"
  "github.com/fengyfei/gu/libs/logger"
  "strings"
  "strconv"
  "time"
  "io/ioutil"
)

func GetNameByTime(path string, suffix string) string {
  files, _ := ioutil.ReadDir(path)
  timeStamp := time.Now().Unix()

  return strconv.FormatInt(timeStamp, 10) + strconv.Itoa(len(files)) + "." + suffix
}

func SavePicture(base64Str string, pathPrefix string) (string, error) {
  slice := strings.Split(base64Str, ",")
  suffix := string([]byte(slice[0])[11:len(slice[0]) - 7]) // picture format suffix

  byteData, err := base64.StdEncoding.DecodeString(slice[1])
  if err != nil {
    logger.Error(err)

    return "", err
  }

  path := "./img/" + pathPrefix
  fileName := GetNameByTime(path, suffix)

  err = ioutil.WriteFile(path + fileName, byteData, 0777)
  if err != nil {
    logger.Error(err)
  }

  return path + fileName, err
}
