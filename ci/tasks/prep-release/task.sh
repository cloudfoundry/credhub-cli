#!/bin/bash

set -eux

export GOPATH=$PWD/go
export GOARCH=amd64
BUILD_ROOT=$PWD

binary_name="cm-cli"

echo ${binary_name} > ${PREP_RELEASE_OUTPUT_PATH}/name
date +'%s' > ${PREP_RELEASE_OUTPUT_PATH}/tag
cd ${GOPATH}/src/github.com/pivotal-cf/cm-cli

make dependencies

GOOS=linux go build -o cm-linux
GOOS=darwin go build -o cm-darwin
GOOS=windows go build -o cm-windows.exe

mv cm-* ${BUILD_ROOT}/${PREP_RELEASE_OUTPUT_PATH}
