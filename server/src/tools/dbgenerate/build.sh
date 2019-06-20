#!/bin/bash
#BASEPATH=$(cd `dirname $0`; pwd)
#cd $BASEPATH
GOPATH=`pwd`
go build -race -ldflags "-extldflags '-static'" -o bin/dbgenerate src/tools/dbgenerate/main.go
