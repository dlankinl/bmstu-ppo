#!/bin/zsh

set -e

SCRIPT_PATH="$(dirname "$(realpath "$0")")"

source "$SCRIPT_PATH/env.bash"

go run main.go