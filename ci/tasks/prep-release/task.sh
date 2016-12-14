#!/bin/bash

set -eux

export GOPATH=$PWD/go
export PATH=$PATH:$GOPATH/bin
export GOARCH=amd64
BUILD_ROOT=$PWD

binary_name="credhub-cli"
build_number=$(date +%s)

echo ${binary_name} > ${PREP_RELEASE_OUTPUT_PATH}/name
echo ${build_number} > ${PREP_RELEASE_OUTPUT_PATH}/tag
cd ${GOPATH}/src/github.com/pivotal-cf/credhub-cli

make dependencies

for os in linux darwin windows; do
  VERSION=${build_number} GOOS=${os} make build
  tar -C build -cvzf ${BUILD_ROOT}/${PREP_RELEASE_OUTPUT_PATH}/"credhub-${os}.tgz" .
  rm -rf build
done
