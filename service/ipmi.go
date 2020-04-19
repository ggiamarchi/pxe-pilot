package service

import (
	"fmt"
	"strings"

	"github.com/ggiamarchi/pxe-pilot/model"

	"github.com/ggiamarchi/pxe-pilot/logger"
)

type IpmiBmcAdapter struct {
}

// PowerStatus is a wrapper for `ipmitool chassis power status`
func (adapter IpmiBmcAdapter) PowerStatus(host *model.Host) (string, error) {
	context := host.IPMI
	stdout, _, err := ipmitool(context, "chassis power status")
	if err != nil || stdout == nil {
		host.PowerState = "Unknown"
		return host.PowerState, err
	}
	logger.Info("ipmitool stdout for IP %s :: %s", context.Hostname, *stdout)
	if strings.Contains(*stdout, "Chassis Power is on") {
		host.PowerState = "On"
		return host.PowerState, nil
	}
	if strings.Contains(*stdout, "Chassis Power is off") {
		host.PowerState = "Off"
		return host.PowerState, nil
	}
	host.PowerState = "Unknown"
	return host.PowerState, nil
}

// PowerOn is a wrapper for `ipmitool chassis power on`
func (adapter IpmiBmcAdapter) PowerOn(host *model.Host) error {
	_, _, err := ipmitool(host.IPMI, "chassis power on")
	return err
}

// PowerReset is a wrapper for `ipmitool chassis power reset`
func (adapter IpmiBmcAdapter) PowerReset(host *model.Host) error {
	_, _, err := ipmitool(host.IPMI, "chassis power reset")
	return err
}

// PowerOff is a wrapper for `ipmitool chassis power off`
func (adapter IpmiBmcAdapter) PowerOff(host *model.Host) error {
	_, _, err := ipmitool(host.IPMI, "chassis power off")
	return err
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

	baseCmd := fmt.Sprintf("ipmitool%s -N 2 -R 2 -H %s -U %s -P %s ", interfaceOpt, context.Hostname, context.Username, context.Password)

	fullCommand := baseCmd + command
	stdout, stderr, err := execCommand(fullCommand)

	return &stdout, &stderr, err
}
