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
 *     Initial: 2017/09/17        Li Zebang
 */

package middleware

import (
	"io/ioutil"
	"net/http"
	"testing"
)

func hello(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("hello"))
}

func TestBasicAuth(t *testing.T) {
	go func() {
		err := http.ListenAndServe(":8080", &BasicAuth{http.HandlerFunc(hello)})
		t.Fatal(err)
	}()

	req, err := http.NewRequest("GET", "http://localhost:8080", nil)
	if err != nil {
		t.Fatal(err)
	}

	req.SetBasicAuth("middleware", "123456")

	c := &http.Client{}
	resp, err := c.Do(req)
	if err != nil {
		t.Fatal(err)
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		t.Fatal(err)
	}

	if string(body) == "" {
		t.Fatal("Authentication Failed!")
	} else {
		t.Log(string(body))
	}
}

func TestSetCookie(t *testing.T) {
	go func() {
		err := http.ListenAndServe(":8081", SetCookie(http.HandlerFunc(hello)))
		t.Fatal(err)
	}()

	req, err := http.NewRequest("GET", "http://localhost:8081", nil)
	if err != nil {
		t.Fatal(err)
	}

	req.SetBasicAuth("middleware", "")

	c := &http.Client{}
	resp, err := c.Do(req)
	if err != nil {
		t.Fatal(err)
	}

	cookies := resp.Cookies()
	if len(cookies) != 0 {
		t.Log(cookies[0].Value)
	} else {
		t.Fatal("Set Cookie Error")
	}
}
