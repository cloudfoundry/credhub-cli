#!/bin/bash

set -exu

PREP_RELEASE_INPUT=$(mktemp -d -t prep-release-input)
echo '{"source":{"interval":"1s"},"version":{"time":"2012-07-18T22:36:54.500564939Z"}}' > $PREP_RELEASE_INPUT/input

PREP_RELEASE_OUTPUT=$(mktemp -d -t prep-release-output)

fly \
  -t private \
  execute \
  -c task.yml \
  -i code-repo=../../.. \
  -i task-repo=../../.. \
  -o prep-release-output=$PREP_RELEASE_OUTPUT
