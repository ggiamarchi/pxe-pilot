#!/bin/bash

set -ex

PXE_PILOT_BIN=pxe-pilot

if [ ! -e ${GOPATH}/src/github.com/ggiamarchi/pxe-pilot/${PXE_PILOT_BIN} ]; then
	PXE_PILOT_BIN=pxe-pilot-linux-amd64
fi

sudo ln -s ${GOPATH}/src/github.com/ggiamarchi/pxe-pilot/${PXE_PILOT_BIN} /usr/local/bin/pxe-pilot

cat ${GOPATH}/src/github.com/ggiamarchi/pxe-pilot/test/pxe-pilot-template.yml |
	sed "s#%%TFTP_ROOT%%#${GOPATH}/src/github.com/ggiamarchi/pxe-pilot/test/tftp_root#" |
	sed "s#%%CONFIG_DIR%%#${GOPATH}/src/github.com/ggiamarchi/pxe-pilot/test/tftp_root/pxelinux.cfg/conf#" \
		> ${GOPATH}/src/github.com/ggiamarchi/pxe-pilot/test/pxe-pilot.yml

sudo bash -c "cat > /etc/systemd/system/pxe-pilot.service" <<- EOF
	[Unit]
	Description=PXE Pilot Server
	After=network-online.target

	[Service]
	User=root
	Group=root
	ExecStart=/usr/local/bin/pxe-pilot server -c ${GOPATH}/src/github.com/ggiamarchi/pxe-pilot/test/pxe-pilot.yml
	KillMode=process
	Restart=on-failure

	[Install]
	WantedBy=multi-user.target
EOF

sudo systemctl enable pxe-pilot
sudo systemctl start pxe-pilot
