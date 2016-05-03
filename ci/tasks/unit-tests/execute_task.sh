#!/bin/bash

set -eu

fly \
  -t private \
  execute \
  -c task.yml \
  -i task-repo=../../../ \
  -i code-repo=../../../
