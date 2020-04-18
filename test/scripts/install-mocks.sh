#!/bin/bash

set -ex

sudo ln -s $GOPATH/src/github.com/ggiamarchi/pxe-pilot/test/mocks/ipmitool-mock /usr/local/bin/
sudo ln -s $GOPATH/src/github.com/ggiamarchi/pxe-pilot/test/mocks/ipmitool-mock /usr/local/bin/ipmitool

sudo mv /usr/sbin/arp /usr/sbin/arp-real
sudo ln -s $GOPATH/src/github.com/ggiamarchi/pxe-pilot/test/mocks/arp-mock /usr/sbin/arp

sudo ln -s $GOPATH/src/github.com/ggiamarchi/pxe-pilot/test/mocks/fping-mock /usr/bin/fping

sudo mkdir -p /etc/pxe-pilot

sudo bash -c "cat > /etc/pxe-pilot/test-home" <<- EOF
	$GOPATH/src/github.com/ggiamarchi/pxe-pilot/test
EOF
