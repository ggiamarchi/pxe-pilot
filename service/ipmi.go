package service

import (
	"bytes"
	"fmt"
	"os/exec"
	"strings"

	"github.com/ggiamarchi/pxe-pilot/model"

	"github.com/ggiamarchi/pxe-pilot/logger"
)

// ChassisPowerStatus is a wrapper for for `ipmitool chassis power status`
func ChassisPowerStatus(context *model.IPMI) (string, error) {
	stdout, _, err := ipmitool(context, "chassis power status")
	if err != nil {
		context.Status = "Unknown"
		return context.Status, err
	}
	if strings.Contains(*stdout, "Chassis Power is on") {
		context.Status = "On"
		return context.Status, nil
	}
	context.Status = "Off"
	return context.Status, nil
}

// ChassisPowerOn is a wrapper for for `ipmitool chassis power on`
func ChassisPowerOn(context *model.IPMI) error {
	_, _, err := ipmitool(context, "chassis power on")
	return err
}

// ChassisPowerReset is a wrapper for for `ipmitool chassis power reset`
func ChassisPowerReset(context *model.IPMI) error {
	_, _, err := ipmitool(context, "chassis power reset")
	return err
}

// ChassisPowerOff is a wrapper for for `ipmitool chassis power off`
func ChassisPowerOff(context *model.IPMI) error {
	_, _, err := ipmitool(context, "chassis power off")
	return err
}

func execCommand(command string, args ...interface{}) (string, string, error) {

	fmtCommand := fmt.Sprintf(command, args...)

	splitCommand := strings.Split(fmtCommand, " ")

	logger.Info("Executing command :: %s :: with args :: %v => %s", command, args, fmtCommand)

	cmdName := splitCommand[0]
	cmdArgs := splitCommand[1:len(splitCommand)]

	cmd := exec.Command(cmdName, cmdArgs...)

	var stdout bytes.Buffer
	cmd.Stdout = &stdout

	var stderr bytes.Buffer
	cmd.Stderr = &stderr

	err := cmd.Run()

	return stdout.String(), stderr.String(), err
}

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

func ipmitool(context *model.IPMI, command string) (*string, *string, error) {

	// Populate IPMI Hostname
	if context.Hostname == "" {
		context.Hostname, _ = getIPFromMAC(context.MACAddress)
	}

	if context.Hostname == "" {
		return nil, nil, logger.Errorf("Unable to find IPMI interface for MAC '%s'", context.MACAddress)
	}

	var interfaceOpt string
	if context.Interface != "" {
		interfaceOpt = fmt.Sprintf(" -I %s", context.Interface)
	}

	baseCmd := fmt.Sprintf("ipmitool%s -N 1 -R 2 -H %s -U %s -P %s ", interfaceOpt, context.Hostname, context.Username, context.Password)

	fullCommand := baseCmd + command
	stdout, stderr, err := execCommand(fullCommand)

	if err != nil {
		logger.Error("IPMI command failed <%s> - %s", fullCommand, err)
	}

	return &stdout, &stderr, err
}
