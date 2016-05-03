#!/bin/bash -eu

export GOPATH=$PWD/go
export GOARCH=amd64
BUILD_ROOT=$PWD
echo "cm-cli" >$PREP_RELEASE_OUTPUT_PATH/name
date +'%s' > $PREP_RELEASE_OUTPUT_PATH/tag
cd $GOPATH/src/github.com/pivotal-cf/cm-cli

for GOOS in linux darwin windows; do
  export GOOS

  go build \
    -a \
    -tags netgo \
    -installsuffix netgo \
    -o $BUILD_ROOT/$PREP_RELEASE_OUTPUT_PATH/cm-${GOOS}
done

