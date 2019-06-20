#!/bin/bash
GOPATH=`pwd`
CGO_ENABLED=0 go build -ldflags "-extldflags '-static'" -o bin/usdtsvr src/usdtsvr/main.go
