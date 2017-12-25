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
 *     Initial: 2017/12/23        Yang Chenglong
 */

package minio

import (
	"testing"

	"fmt"
	"github.com/minio/minio-go"
)

var (
	c   *Client
	err error
)

func TestNewClient(t *testing.T) {
	config := Config{
		endpoint:        "127.0.0.1:9001",
		accessKeyID:     "minio",
		secretAccessKey: "minio20171222",
		useSSL:          false,
	}
	c, err = NewClient(&config)
	if err != nil {
		t.Fatal(err)
	}
}

func TestClient_NewBucket(t *testing.T) {
	bucketName := "myfile"
	location := "us-east-8"

	err = c.NewBucket(bucketName, location)
	if err != nil {
		t.Fatal(err)
	}
}

func TestClient_PutObject(t *testing.T) {
	bucketName := "myfile"
	objectName := "test.zip"
	filePath := "./test.zip"
	contentType := "application/zip"

	// Upload the zip file with FPutObject
	n, err := c.PutObject(bucketName, objectName, filePath, minio.PutObjectOptions{ContentType: contentType})
	if err != nil {
		t.Fatal(err)
	}

	t.Log("the size of file is:", n)
}

func TestClient_PutImage(t *testing.T) {
	bucketName := "myfile"
	objectName := "test.png"
	filePath := "./test.png"
	n, err := c.PutImage(bucketName, objectName, filePath)
	if err != nil {
		t.Fatal(err)
	}

	t.Log("the size of file is:", n)
}

func TestClient_PutText(t *testing.T) {
	bucketName := "myfile"
	objectName := "T_test.txt"
	filePath := "./Ttest.txt"
	n, err := c.PutImage(bucketName, objectName, filePath)
	if err != nil {
		t.Fatal(err)
	}

	t.Log("the size of file is:", n)
}

func TestClient_ListBuckets(t *testing.T) {
	buckets, err := c.ListBuckets()
	if err != nil {
		t.Fatal(err)
	}
	for _, bucket := range buckets {
		fmt.Println("bucket:", bucket.Name)
	}

}

func TestClient_ListObjects(t *testing.T) {
	doneCh := make(chan struct{})
	defer close(doneCh)
	info := c.ListObjects("myfile", "", true, doneCh)
	for v := range info {
		if v.Err != nil {
			t.Log(v.Err)
			return
		}
		t.Log(v)
	}
}
func TestClient_GetFile(t *testing.T) {
	bucketName := "myfile"
	objectName := "test.zip"
	filePath := "./test1.zip"

	err := c.GetObject(bucketName, objectName, filePath, minio.GetObjectOptions{})
	if err != nil {
		t.Fatal(err)
	}
}

func TestClient_DeleteFie(t *testing.T) {
	bucketName := "myfile"
	objectName := "test.zip"
	err := c.DeleteObject(bucketName, objectName)
	if err != nil {
		t.Fatal(err)
	}
}
