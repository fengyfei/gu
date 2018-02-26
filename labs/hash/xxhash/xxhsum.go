// Copyright (c) 2016 Caleb Spare

// MIT License

// Permission is hereby granted, free of charge, to any person obtaining
// a copy of this software and associated documentation files (the
// "Software"), to deal in the Software without restriction, including
// without limitation the rights to use, copy, modify, merge, publish,
// distribute, sublicense, and/or sell copies of the Software, and to
// permit persons to whom the Software is furnished to do so, subject to
// the following conditions:

// The above copyright notice and this permission notice shall be
// included in all copies or substantial portions of the Software.

// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND,
// EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF
// MERCHANTABILITY, FITNESS FOR A PARTICULAR PURPOSE AND
// NONINFRINGEMENT. IN NO EVENT SHALL THE AUTHORS OR COPYRIGHT HOLDERS BE
// LIABLE FOR ANY CLAIM, DAMAGES OR OTHER LIABILITY, WHETHER IN AN ACTION
// OF CONTRACT, TORT OR OTHERWISE, ARISING FROM, OUT OF OR IN CONNECTION
// WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE SOFTWARE.

package main

import (
	"fmt"
	"io"
	"os"

	"github.com/cespare/xxhash"
)

func main() {
	if contains(os.Args[1:], "-h") {
		fmt.Fprintf(os.Stderr, `Usage:
  %s [filenames]
If no filenames are provided or only - is given, input is read from stdin.
`, os.Args[0])
		os.Exit(1)
	}
	if len(os.Args) < 2 || len(os.Args) == 2 && string(os.Args[1]) == "-" {
		printHash(os.Stdin, "-")
		return
	}
	for _, path := range os.Args[1:] {
		f, err := os.Open(path)
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			continue
		}
		printHash(f, path)
		f.Close()
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

func printHash(r io.Reader, name string) {
	h := xxhash.New()
	if _, err := io.Copy(h, r); err != nil {
		fmt.Fprintln(os.Stderr, err)
		return
	}
	fmt.Printf("%016x  %s\n", h.Sum64(), name)
}
