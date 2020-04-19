#!/bin/bash

set -ex

sudo rm -rf /usr/local/go $HOME/gopath
wget https://dl.google.com/go/go1.14.linux-amd64.tar.gz 2> /dev/null
sudo tar -C /usr/local -xzf go1.14.linux-amd64.tar.gz
mkdir $HOME/gopath
echo 'export GOROOT=/usr/local/go' >> $HOME/.profile
echo 'export GOPATH=$HOME/gopath' >> $HOME/.profile
echo 'export PATH=$PATH:$GOROOT/bin:$GOPATH/bin' >> $HOME/.profile

source $HOME/.profile

mkdir -p $GOPATH/src/github.com/ggiamarchi
ln -s /vagrant $GOPATH/src/github.com/ggiamarchi/pxe-pilot
