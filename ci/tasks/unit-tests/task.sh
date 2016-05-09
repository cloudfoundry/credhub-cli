#!/bin/bash

set -exu

export GOPATH=$PWD/go
export PATH=$PATH:$GOPATH/bin

cd go/src/github.com/pivotal-cf/cm-cli
make test
