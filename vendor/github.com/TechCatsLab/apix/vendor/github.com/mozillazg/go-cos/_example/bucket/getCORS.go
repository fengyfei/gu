package main

import (
	"context"
	"fmt"
	"net/url"
	"os"

	"net/http"

	"github.com/mozillazg/go-cos"
	"github.com/mozillazg/go-cos/debug"
)

func main() {
	u, _ := url.Parse(os.Getenv("COS_BUCKET_URL"))
	b := &cos.BaseURL{
		BucketURL: u,
	}
	c := cos.NewClient(b, &http.Client{
		Transport: &cos.AuthorizationTransport{
			SecretID:  os.Getenv("COS_SECRETID"),
			SecretKey: os.Getenv("COS_SECRETKEY"),
			Transport: &debug.DebugRequestTransport{
				RequestHeader:  true,
				RequestBody:    true,
				ResponseHeader: true,
				ResponseBody:   true,
			},
		},
	})

	v, _, err := c.Bucket.GetCORS(context.Background())
	if err != nil {
		panic(err)
	}
	for _, r := range v.Rules {

		fmt.Printf("%s, %s\n", r.AllowedOrigins, r.AllowedMethods)
	}
}