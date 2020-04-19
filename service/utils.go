package service

import (
	"bytes"
	"fmt"
	"os/exec"
	"strings"

	"github.com/ggiamarchi/pxe-pilot/logger"
)

// getIPFromMAC reads the ARP table to find the IP address matching the given MAC address
func getIPFromMAC(mac string) (string, error) {

	stdout, _, err := execCommand("sudo arp -an")

	if err != nil {
		return "", err
	}

	lines := strings.Split(stdout, "\n")

	for _, v := range lines {
		if strings.TrimSpace(v) == "" {
			continue
		}
		fields := strings.Fields(v)

		if normalizeMACAddress(mac) == normalizeMACAddress(fields[3]) {
			return fields[1][1 : len(fields[1])-1], nil
		}
	}

	return "", nil
}

// normalizeMACAddress takes the input MAC address and remove every non hexa symbol
// and lowercase everything else
func normalizeMACAddress(mac string) string {
	var buffer bytes.Buffer

	macArray := strings.Split(strings.ToLower(mac), ":")

	for i := 0; i < len(macArray); i++ {
		m := macArray[i]
		if len(m) == 1 {
			buffer.WriteByte(byte('0'))
		}
		for j := 0; j < len(m); j++ {
			if isHexChar(m[j]) {
				buffer.WriteByte(m[j])
			}
		}
	}
	return buffer.String()
}

func isHexChar(char byte) bool {
	switch char {
	case
		byte('a'), byte('b'), byte('c'), byte('d'),
		byte('e'), byte('f'), byte('0'), byte('1'),
		byte('2'), byte('3'), byte('4'), byte('5'),
		byte('6'), byte('7'), byte('8'), byte('9'):
		return true
	}
	return false
}

func execCommand(command string, args ...interface{}) (string, string, error) {

	fmtCommand := fmt.Sprintf(command, args...)

	splitCommand := strings.Split(fmtCommand, " ")

	logger.Info("Executing command :: %s :: with args :: %v => %s", command, args, fmtCommand)

	cmdName := splitCommand[0]
	cmdArgs := splitCommand[1:]

	cmd := exec.Command(cmdName, cmdArgs...)

	var stdout bytes.Buffer
	cmd.Stdout = &stdout

	var stderr bytes.Buffer
	cmd.Stderr = &stderr

	err := cmd.Run()

	if err != nil {
		logger.Error("Command failed <%s> | error : %s | stdout : %s | stderr : %s", fmtCommand, err, stdout.String(), stderr.String())
	}

	return stdout.String(), stderr.String(), err
}
