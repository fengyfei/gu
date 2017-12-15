#! /bin/sh

linuxRelease=admin.linux

GOOS=linux GOARCH=amd64 go build -o $linuxRelease
