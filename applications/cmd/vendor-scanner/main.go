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
)

const (
	vendor = "vendor"
)

func main() {
	if sliceContains(os.Args[1:], "-h") || len(os.Args) > 2 {
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
			vendorPath, _ = filepath.Abs(vendorPath)
			fmt.Printf("---%s---\n", vendorPath[:len(vendorPath)-7])

			var paths []string
			err := filepath.Walk(vendorPath, func(path string, info os.FileInfo, err error) error {
				if err != nil {
					return err
				}

				depth := strings.Count(path, "/") - strings.Count(vendorPath, "/") - 1
				siteDepth := mapContains(sitesDepth, path)
				if depth == siteDepth {
					if info.IsDir() {
						paths = append(paths, path)
						fmt.Println(strings.SplitN(path, "/"+vendor+"/", 2)[1])
					} else {
						path = path[:strings.LastIndex(path, "/")]
						if !sliceContains(paths, path) {
							paths = append(paths, path)
							fmt.Println(strings.SplitN(path, "/"+vendor+"/", 2)[1])
						}
					}
				}

				return nil
			})

			if err != nil {
				fmt.Fprintf(os.Stderr, `Usage:
  %s [pathname]
Error in walking the file tree : %s
`, os.Args[0], err)
				os.Exit(1)
			}

			return nil
		}
		return nil
	})

	if err != nil {
		fmt.Fprintf(os.Stderr, `Usage:
  %s [pathname]
Error in walking the file tree : %s
`, os.Args[0], err)
		os.Exit(1)
	}
}

func sliceContains(slice []string, str string) bool {
	for _, s := range slice {
		if str == s {
			return true
		}
	}
	return false
}

func mapContains(m map[string]int, str string) int {
	str = strings.Split(str[strings.Index(str, "/"+vendor+"/")+8:], "/")[0]
	for k, v := range m {
		if str == k {
			return v
		}
	}
	return defaultDepth
}
