#!/bin/bash

set -eux

export GOPATH=$PWD/go
export PATH=$PATH:$GOPATH/bin

cd go/src/github.com/pivotal-cf/credhub-cli
make ci
