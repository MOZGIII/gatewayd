#!/bin/bash
DIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd )"
PROJECT="$(basename "$DIR")"

mkdir -p "$DIR/.goworkspace/src"
ln -sTf "../.." "$DIR/.goworkspace/src/$PROJECT"
