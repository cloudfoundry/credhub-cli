#!/usr/bin/env bash

function set_bash_error_handling() {
    set -euo pipefail
}

function go_to_project_root_directory() {
    local -r script_dir=$( dirname "${BASH_SOURCE[0]}")

    cd "$script_dir/.."
}

function check_ssh_key() {
    if ! ssh-add -l >/dev/null; then
        echo "No SSH key loaded! Please run vkl."
        exit 1
    fi
}

function run_linters() {
    ./scripts/lint.sh
}

function login_to_local_credhub(){
    echo "Logging in to CredHub"
    credhub a https://localhost:9000 --skip-tls-validation
    credhub l -u credhub -p password
}

function start_background_server(){
    pushd ~/workspace/credhub-release/src/credhub
    ./scripts/start_server.sh &> /dev/null &
    export PID=$!
    echo "Waiting for server to start up"
    sleep 30

    COUNTER=1
    while ! curl -s https://localhost:9000/health --insecure > /dev/null; do
        echo "Connection attempt #" ${COUNTER}
        COUNTER=$((COUNTER + 1))
        sleep 10
      done;
    popd
}

function check_for_local_server(){
    echo "Checking for locally running server"
    if curl -s https://localhost:9000/health --insecure > /dev/null; then
      echo "Found locally running CredHub"
    else
      echo "CredHub is not running, attempting to start"
      start_background_server
    fi
    login_to_local_credhub
}


function run_tests() {
    export GOPATH=~/go
    pushd ${GOPATH}/src/github.com/cloudfoundry-incubator/credhub-acceptance-tests
        ./scripts/run_tests.sh
    popd
}

function fail_for_uncommitted_changes() {
    local -r number_of_uncommitted_changes=$(git status -s | wc -l | tr -d '[:space:]')

    if [ "$number_of_uncommitted_changes" != "0" ]; then
        echo "WARNING: uncommitted changes detected"
        exit 1
    fi
}

function push_code() {
    git push
}

function display_ascii_success_message() {
    local -r green_color_code='\033[1;32m'
    echo -e "${green_color_code}\\n$(cat scripts/success_ascii_art.txt)"
}

function main() {
    set_bash_error_handling
    go_to_project_root_directory
    check_ssh_key
    check_for_local_server

    run_linters
    fail_for_uncommitted_changes

    run_tests

    push_code
    display_ascii_success_message
    kill ${PID}
}

main
