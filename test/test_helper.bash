#!/bin/bash

BATS_OUT_LOG=${BATS_TEST_DIRNAME}/${TEST_LOG_FILE:-test.log}

setup() {
	{
		title "#${BATS_TEST_NUMBER} | Begin test | ${BATS_TEST_NAME} | ${BATS_TEST_FILENAME}"

		TEST_HOME=$(cat /etc/pxe-pilot/test-home)

		#
		# Emulate ARP table flushing
		#
		sudo rm -f /tmp/arp-populated

		#
		# Delete emulated IPMI power status (all to power off)
		#
		find ${TEST_HOME}/mocks/data -name status | xargs rm -f

		#
		# Delete PXE Pilot deployed configurations
		#
		TFTP_ROOT=${TEST_HOME}/tftp_root
		rm -f ${TFTP_ROOT}/pxelinux.cfg/01-*

		#
		# '(re)start PXE Piolt server'
		#
		sleep 0.5
		sudo systemctl restart pxe-pilot
		sleep 0.5

	} >> ${BATS_OUT_LOG}
}

teardown() {
	flush_log
	{
		title "#${BATS_TEST_NUMBER} | End test   | ${BATS_TEST_NAME} | ${BATS_TEST_FILENAME}"
	} >> ${BATS_OUT_LOG}
}

title() {
	{
		echo ""
		echo "---------------------------------------------------------------------------------------------------------------------"
		echo "--- ${1}"
		echo "---------------------------------------------------------------------------------------------------------------------"
		echo ""
	} >> ${BATS_OUT_LOG}
}

flush_log() {
	{
		echo "${lines[@]}"
	} >> ${BATS_OUT_LOG}
}
