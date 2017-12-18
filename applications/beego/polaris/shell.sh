#! /bin/sh

linuxRelease=polaris.linux

GOOS=linux GOARCH=amd64 go build -o $linuxRelease
