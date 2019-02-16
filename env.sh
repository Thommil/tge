#!/bin/bash

if [ -z "$GOPATH" ] ; then
    export GOPATH="$( cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd )/go"
    export PATH=$PATH:$GOPATH/bin
fi