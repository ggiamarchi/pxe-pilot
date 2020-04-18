#!/bin/bash

set -ex

cd $GOPATH/src/github.com/ggiamarchi/pxe-pilot
rm -rf bin
make dep-dev
make
