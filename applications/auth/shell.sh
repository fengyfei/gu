#! /bin/sh

linuxRelease=auth.linux

GOOS=linux GOARCH=amd64 go build -o $linuxRelease
