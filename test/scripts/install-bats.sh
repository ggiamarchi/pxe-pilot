#!/bin/bash

set -ex

pushd /tmp/
git clone https://github.com/bats-core/bats-core.git
pushd bats-core
sudo ./install.sh /usr/local
popd
popd
