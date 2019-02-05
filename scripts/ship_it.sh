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
    go fmt $(go list ./... | grep -v /vendor/)
}

function run_tests() {
    make test
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

    run_linters
    fail_for_uncommitted_changes

    run_tests

    push_code
    display_ascii_success_message
}

main
