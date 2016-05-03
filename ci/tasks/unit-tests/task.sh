#!/bin/bash

set -exu

export GOPATH=$PWD/go
export PATH=$PATH:$GOPATH/bin

cd go/src/github.com/pivotal-cf/cm-cli
go get github.com/tools/godep
godep restore ./...
go install github.com/onsi/ginkgo/ginkgo
ginkgo -r
