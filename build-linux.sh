#!/bin/bash
export GOPATH=`pwd`
export GOARCH=amd64
export GOOS=linux
cd bin

go build -o server -ldflags "-X main.VERSION=1.0.4 -X 'main.BUILD_TIME=`date`' -s -w" ../src/server.go
go build -o robot -ldflags "-w -s" ../src/robot.go

upx -9 server
upx -9 robot

read -p "Press any key to continue." var


