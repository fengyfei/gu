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
 *     Initial: 2018/02/25        Feng Yifei
 *     Modify:  2018/02/25        Li Zebang
 */

package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/fengyfei/gu/libs/logger"
)

const (
	vendor = "vendor"
)

func main() {
	if contains(os.Args[1:], "-h") || len(os.Args) > 2 {
		fmt.Fprintf(os.Stderr, `Usage:
  %s [pathname]
`, os.Args[0])
		os.Exit(1)
	}

	var dir string
	if len(os.Args) == 1 {
		dir = "."
	} else {
		dir = os.Args[1]
	}

	err := filepath.Walk(dir, func(vendorPath string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if info.IsDir() && info.Name() == vendor {
			fmt.Printf("----------%s----------\n", vendorPath)
			err := filepath.Walk(vendorPath, func(path string, info os.FileInfo, err error) error {
				if err != nil {
					return err
				}

				for site, siteDepth := range sitesDepth {
					depth := strings.Count(path, "/") - strings.Count(vendorPath, "/") - 1
					if strings.Contains(path, site) {
						if depth == siteDepth {
							fmt.Println(path)
							return nil
						}
					} else {
						if depth == defaultDepth {
							fmt.Println(path)
							return nil
						}
					}
				}

				return nil
			})

			if err != nil {
				logger.Error("error in walking the file tree rooted at", vendorPath, err)
				return err
			}

			return nil
		}
		return nil
	})

	if err != nil {
		logger.Error("error in walking the file tree rooted at", dir, err)
	}
}

func contains(ss []string, s string) bool {
	for _, s1 := range ss {
		if s1 == s {
			return true
		}
	}
	return false
}
