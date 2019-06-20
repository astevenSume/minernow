#!/bin/bash
GOPATH=`pwd`
#go build -race -ldflags "-extldflags '-static'" -o bin/gatewaysvr src/gatewaysvr/main.go
CGO_ENABLED=0 go build -ldflags "-extldflags '-static'" -o bin/admin src/admin/main.go
