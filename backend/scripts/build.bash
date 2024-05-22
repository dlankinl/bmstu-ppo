#!/bin/zsh
scriptPath=$(dirname $(realpath $0))
cd $scriptPath/../src
rice clean
go build -o $scriptPath/../bin/server