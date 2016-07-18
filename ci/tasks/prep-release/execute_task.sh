#!/bin/bash

set -exu

PREP_RELEASE_INPUT=$(mktemp -d -t prep-release-input)
echo "$(date +%s)" > $PREP_RELEASE_INPUT/input

PREP_RELEASE_OUTPUT=$(mktemp -d -t prep-release-output)

fly \
  -t private \
  execute \
  -c task.yml \
  -i clock=$PREP_RELEASE_INPUT \
  -i code-repo=../../.. \
  -i task-repo=../../.. \
  -o prep-release-output=$PREP_RELEASE_OUTPUT
