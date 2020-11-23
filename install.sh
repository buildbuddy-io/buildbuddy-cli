#!/bin/bash

read -r -d '' USAGE << EOM
Usage: $0 [-b]
       -b   -- If set, installs beta version of $0 CLI.
EOM

function usage() {
    echo "$USAGE"
    exit 1
}

function old_install() {
    curl -fsSL -o bb https://raw.githubusercontent.com/buildbuddy-io/cli/master/bb && chmod 755 bb && sudo mv bb /usr/local/bin/bb
    exit 0
}

function new_install() {
    arch=`uname -m`  # x86_64 | i386
    if [ "$arch" == "x86_64" ]; then
	arch="amd64"
    fi
    os=`uname -s` # Linux | Darwin
    os="$(tr [A-Z] [a-z] <<< "$os")"
    curl -fsSL -o buildbuddy https://github.com/buildbuddy-io/cli/releases/latest/download/buildbuddy-$os-$arch && chmod 0755 buildbuddy && sudo mv buildbuddy /usr/local/bin/buildbuddy
    exit 0
}

# new world: detect arch and install go tool
# old world: install the python script

while getopts ":b" arg; do
    case $arg in
	b)
	    new_install
	    ;;
	*)
	    usage
	    ;;
    esac
done
if [ "$arg" == "?" ]; then
    old_install
fi
