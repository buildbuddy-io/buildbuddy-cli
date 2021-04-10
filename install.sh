#!/usr/bin/env bash

function new_install() {
    arch=`uname -m`  # x86_64 | i386
    if [ "$arch" == "x86_64" ]; then
	arch="amd64"
    fi
    os=`uname -s` # Linux | Darwin
    os="$(tr [A-Z] [a-z] <<< "$os")"
    tmpFile=$(mktemp buildbuddy.XXXXX)
    trap "rm -f $tmpFile" 0 2 3 15
    curl -fsSL -o $tmpFile https://github.com/buildbuddy-io/cli/releases/latest/download/buildbuddy-$os-$arch && chmod 0755 $tmpFile && sudo mv $tmpFile /usr/local/bin/bb
    exit 0
}

new_install