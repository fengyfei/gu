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

package main

import (
	"github.com/fengyfei/gu/libs/logger"
	"github.com/minio/minio-go"
)

// Config contains the necessary configuration about Client.
type Config struct {
	endpoint        string
	accessKeyID     string
	secretAccessKey string
	useSSL          bool
}

// Client represents a minio client.
type Client struct {
	client *minio.Client
	config *Config
}

// NewClient creates a new minio client.
func NewClient(c *Config) (*Client, error) {
	var (
		err    error
		client Client
	)

	client.client, err = minio.New(c.endpoint, c.accessKeyID, c.secretAccessKey, c.useSSL)
	if err != nil {
		return nil, err
	}

	client.config = c

	return &client, nil
}

// NewBucket creates a new bucket.
func (c *Client) NewBucket(name string, timezone string) error {
	err := c.client.MakeBucket(name, timezone)
	if err != nil {
		// Check to see if we already own this bucket (which happens if you run this twice)
		exists, err := c.client.BucketExists(name)
		if err == nil && exists {
			logger.Warn("We already own", name)
			return nil
		} else {
			logger.Error("make bucket error:", err)
			return err
		}
	}
	return nil
}

// PutFile creates an object in a bucket, with contents from file at filePath.
func (c *Client) PutFile(bucketName, objectName, filePath string, opts minio.PutObjectOptions) (int64, error) {
	n, err := c.client.FPutObject(bucketName, objectName, filePath, opts)
	if err != nil {
		logger.Error("Put file error:", err)
		return 0, err
	}

	return n, nil
}

// GetFile downloads contents of an object to a local file.
func (c *Client) GetFile(bucketName, objectName, filePath string, opts minio.GetObjectOptions) error {
	return c.client.FGetObject(bucketName, objectName, filePath, opts)
}

// DeleteFie removes an object from a bucket.
func (c *Client) DeleteFie(bucketName, objectName string) error {
	return c.client.RemoveObject(bucketName, objectName)
}
