#!/bin/bash

echo "build linux amd64"
GOOS=linux GOARCH=amd64 go build
tar czvf linux_amd64.tar.gz  watch

echo "build linux i386"
GOOS=linux GOARCH=386 go build 
tar czvf linux_i386.tar.gz  watch

echo "build windows amd64"
GOOS=windows GOARCH=amd64 go build -o watch.exe
tar czvf windows_amd64.tar.gz  watch.exe

echo "build windows i386"
GOOS=windows GOARCH=386 go build -o watch.exe
tar czvf windows_i386.tar.gz  watch.exe

echo "build Mac OS X 64bit"
GOOS=darwin GOARCH=amd64 go build
tar czvf darwin_amd64.tar.gz  watch