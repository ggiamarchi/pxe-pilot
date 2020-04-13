package service

import (
	"bytes"
	"regexp"
	"text/template"

	"github.com/ggiamarchi/pxe-pilot/logger"
	"github.com/ggiamarchi/pxe-pilot/model"
)

type HostManagementAdapter interface {
	PowerStatus(host *model.Host) (string, error)
	PowerOn(host *model.Host) error
	PowerReset(host *model.Host) error
	PowerOff(host *model.Host) error
}

type GenericHostManagementAdapter struct {
}

func (adapter GenericHostManagementAdapter) executeManagementCommand(host *model.Host, hostManagementCmd *model.HostManagementCmd) (string, error) {

	var err error

	rm := host.Management

	// Populate host management IP

	if rm.IPAddress == "" {
		rm.IPAddress, _ = getIPFromMAC(rm.MacAddress)
	}

	if rm.IPAddress == "" {
		logger.Error("Unable to find host management IP address for host %s", host.MACAddresses[0])
		return "", nil
	}

	logger.Debug("Host management IP address %s found for host %s", rm.IPAddress, host.MACAddresses[0])

	// Create model for template generation

	model := make(map[string]string)
	model["ip_address"] = rm.IPAddress

	for _, cmdVar := range host.Management.Vars {
		model[cmdVar.Name] = cmdVar.Value
	}

	for _, cmdVar := range host.Management.Adapter.Vars {
		tpl, err := template.New("status").Parse(cmdVar.Value)

		if err != nil {
			return "", err
		}

		var varValue bytes.Buffer

		err = tpl.Execute(&varValue, model)

		if err != nil {
			return "", err
		}

		model[cmdVar.Name] = varValue.String()
	}

	// Compile template

	tpl, err := template.New("status").Parse(hostManagementCmd.Cmd)

	if err != nil {
		return "", err
	}

	// Execute template

	var cmd bytes.Buffer

	err = tpl.Execute(&cmd, model)

	if err != nil {
		return "", err
	}

	// Execute system command

	stdout, _, err := execCommand(cmd.String())

	return stdout, err
}

func (adapter GenericHostManagementAdapter) PowerStatus(host *model.Host) (string, error) {
	logger.Info("Requesting power status for host %s", host.MACAddresses[0])

	powerStatusInfo := host.Management.Adapter.PowerStatus
	stdout, err := adapter.executeManagementCommand(host, powerStatusInfo)

	logger.Debug("Power status command stdout for host %s :: mgmt IP %s :: %s", host.MACAddresses[0], host.Management.IPAddress, stdout)

	if err != nil {
		host.PowerState = "Unknown"
		return host.PowerState, err
	}

	// Check host management command output

	match, err := regexp.MatchString(powerStatusInfo.PatternOn, stdout)
	if err != nil {
		logger.Error("%s", err)
		host.PowerState = "Unknown"
		return host.PowerState, err
	}
	if match {
		host.PowerState = "On"
		return host.PowerState, nil
	}

	match, err = regexp.MatchString(powerStatusInfo.PatternOff, stdout)
	if err != nil {
		logger.Error("%s", err)
		host.PowerState = "Unknown"
		return host.PowerState, err
	}
	if match {
		host.PowerState = "Off"
		return host.PowerState, nil
	}

	host.PowerState = "Unknown"
	return host.PowerState, nil
}

func (adapter GenericHostManagementAdapter) PowerOn(host *model.Host) error {
	logger.Info("Requesting power on for host %s", host.MACAddresses[0])
	stdout, err := adapter.executeManagementCommand(host, host.Management.Adapter.PowerOn)
	logger.Debug("Power on command stdout for host %s :: mgmt IP %s :: %s", host.MACAddresses[0], host.Management.IPAddress, stdout)
	return err
}

func (adapter GenericHostManagementAdapter) PowerReset(host *model.Host) error {
	logger.Info("Requesting power reset for host %s", host.MACAddresses[0])
	stdout, err := adapter.executeManagementCommand(host, host.Management.Adapter.PowerReset)
	logger.Debug("Power reset command stdout for host %s :: mgmt IP %s :: %s", host.MACAddresses[0], host.Management.IPAddress, stdout)
	return err
}

func (adapter GenericHostManagementAdapter) PowerOff(host *model.Host) error {
	logger.Info("Requesting power off for host %s", host.MACAddresses[0])
	stdout, err := adapter.executeManagementCommand(host, host.Management.Adapter.PowerOff)
	logger.Debug("Power on command stdout for host %s :: mgmt IP %s :: %s", host.MACAddresses[0], host.Management.IPAddress, stdout)
	return err
}
