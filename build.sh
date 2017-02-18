#!/bin/bash

echo "build linux amd64"
GOOS=linux GOARCH=amd64 go build -o watch_amd64
tar czvf linux_amd64.tar.gz  watch_amd64

echo "build linux i386"
GOOS=linux GOARCH=386 go build -o watch_i386
tar czvf linux_i386.tar.gz  watch_i386

echo "build windows amd64"
GOOS=windows GOARCH=amd64 go build -o watch_amd64.exe
tar czvf windows_amd64.tar.gz  watch_amd64.exe

echo "build windows i386"
GOOS=windows GOARCH=386 go build -o watch_i386.exe
tar czvf windows_i386.tar.gz  watch_i386.exe

echo "build Mac OS X 64bit"
GOOS=darwin GOARCH=amd64 go build