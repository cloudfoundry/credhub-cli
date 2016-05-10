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

for os in linux darwin windows; do
  GOOS=${os} make build
  tar -C build -cvzf ${BUILD_ROOT}/${PREP_RELEASE_OUTPUT_PATH}/"cm-${os}.tgz" .
  rm -rf build
done
