#!/bin/bash
GOPATH=`pwd`
go build -race -ldflags "-extldflags '-static'" -o bin/errorgenerate src/tools/errorgenerate/main.go
