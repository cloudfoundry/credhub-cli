#!/bin/bash

set -exu

PREP_RELEASE_OUTPUT=$(mktemp -d -t prep-release-output)

fly \
  -t private \
  execute \
  -c task.yml \
  -i code-repo=../../.. \
  -i task-repo=../../.. \
  -o prep-release-output=$PREP_RELEASE_OUTPUT

# rm -rf $PREP_RELEASE_OUTPUT
