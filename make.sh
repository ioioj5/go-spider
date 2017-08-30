#!/bin/bash
pwd=`pwd`

export GOPATH=$pwd

go get github.com/go-sql-driver/mysql

go build -v -o bin/parse parse
