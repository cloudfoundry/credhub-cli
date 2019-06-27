#!/usr/bin/env bash

function set_bash_error_handling() {
    set -euo pipefail
}

function go_to_project_root_directory() {
    local -r script_dir=$( dirname "${BASH_SOURCE[0]}")

    cd "$script_dir/.."
}

function lint_scripts() {
    shellcheck scripts/*.sh
}

function lint_go_code() {
    # shellcheck disable=SC2046
    go fmt $(go list ./... | grep -v /vendor/)
}

function main() {
    set_bash_error_handling
    go_to_project_root_directory
    lint_scripts
    lint_go_code
}

main
