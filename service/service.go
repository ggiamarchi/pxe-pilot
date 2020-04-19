package service

import (
	"io/ioutil"
	"os"
	"sort"
	"strings"
	"sync"

	"fmt"

	"github.com/ggiamarchi/pxe-pilot/common/utils"
	"github.com/ggiamarchi/pxe-pilot/logger"
	"github.com/ggiamarchi/pxe-pilot/model"
)

func Refresh(appConfig *model.AppConfig) error {

	hosts := ReadHosts(appConfig, false)

	m := make(map[string]bool)

	for _, h := range hosts {
		if h.IPMI != nil {
			m[h.IPMI.Subnet] = true
			// Remove IPMI IP to ensure it is refreshed on next IPMI call
			h.IPMI.Hostname = ""
		}
		if h.Management != nil {
			m[h.Management.Subnet] = true
			// Remove IPMI IP to ensure it is refreshed on next IPMI call
			h.Management.IPAddress = ""
		}
	}
	subnets := make([]string, 0, len(m))
	for cidr := range m {
		subnets = append(subnets, cidr)
	}

	logger.Info("Discovery on subnets %+v", subnets)

	var wg sync.WaitGroup
	for _, cidr := range subnets {
		wg.Add(1)
		go func(cidr string) {
			defer wg.Done()
			_, _, err := execCommand("fping -c 1 -D -q -g %s", cidr)
			if err != nil {
				logger.Error("%s", err)
			}
		}(cidr)
	}

	wg.Wait()
	return nil
}

func ReadConfigurations(appConfig *model.AppConfig) []*model.Configuration {
	configurations := make([]*model.Configuration, 0)
	for _, bootloader := range appConfig.Configuration.Bootloaders {
		files, _ := ioutil.ReadDir(appConfig.Configuration.Directory + "/" + bootloader.Name)

		// Sort by age, keeping original order or equal elements.
		sort.SliceStable(files, func(i, j int) bool {
			return strings.Compare(files[i].Name(), files[j].Name()) < 1
		})

		for _, f := range files {
			configuration := &model.Configuration{
				Name:       f.Name(),
				Bootloader: bootloader,
			}
			configurations = append(configurations, configuration)
		}
	}
	return configurations
}

func ReadBootloaders(appConfig *model.AppConfig) []*model.Bootloader {
	return appConfig.Configuration.Bootloaders
}

func OnHost(appConfig *model.AppConfig, host *model.Host) error {
	management, err := getHostManagementAdapter(appConfig, host)
	if err != nil {
		return err
	}
	if management == nil {
		return logger.Errorf("Host %s has no host management interface", host.MACAddresses[0])
	}
	return management.PowerOn(host)
}

func OffHost(appConfig *model.AppConfig, host *model.Host) error {
	management, err := getHostManagementAdapter(appConfig, host)
	if err != nil {
		return err
	}
	if management == nil {
		return logger.Errorf("Host %s has no host management interface", host.MACAddresses[0])
	}
	return management.PowerOff(host)
}

func RebootHost(appConfig *model.AppConfig, host *model.Host) error {
	management, err := getHostManagementAdapter(appConfig, host)
	if err != nil {
		return err
	}
	if management == nil {
		return logger.Errorf("Host %s has no host management interface", host.MACAddresses[0])
	}
	switch status, err := management.PowerStatus(host); status {
	case "On":
		return management.PowerReset(host)
	case "Off":
		return management.PowerOn(host)
	case "Unknown":
		return err
	}
	return logger.Errorf("Reboot host '%s' : Unknown error", host.Name)
}

func ReadHosts(appConfig *model.AppConfig, status bool) []*model.Host {

	pxelinuxDir := appConfig.Tftp.Root + "/pxelinux.cfg"

	var wg sync.WaitGroup

	for i, host := range appConfig.Hosts {

		if status {
			management, _ := getHostManagementAdapter(appConfig, host)
			logger.Info("Host management adapter for host %s :: %T", host.MACAddresses[0], management)

			if management != nil {
				wg.Add(1)
				go func(adapter HostManagementAdapter, hostlocal *model.Host) {
					defer wg.Done()
					_, err := adapter.PowerStatus(hostlocal)
					if err != nil {
						// Retry once
						_, err = adapter.PowerStatus(hostlocal)
						if err != nil {
							logger.Error("Unable to find power status for host %s", hostlocal.MACAddresses[0])
						}
					}
				}(management, host)
			}
		}

		pxeFile := utils.PXEFilenameFromMAC(appConfig.Hosts[i].MACAddresses[0])
		pxeFilePath := fmt.Sprintf("%s/%s", pxelinuxDir, pxeFile)

		if _, err := os.Stat(pxeFilePath); err != nil {
			// No PXE config deployed for this host
			continue
		}

		configFile, err := os.Readlink(pxeFilePath)
		if err != nil {
			panic(err)
		}

		for _, c := range ReadConfigurations(appConfig) {
			if c.Name == configFile[strings.LastIndex(configFile, "/")+1:] {
				appConfig.Hosts[i].Configuration = c
				break
			}
		}
	}

	wg.Wait()

	return appConfig.Hosts
}

type PXEError struct {
	Msg  string
	Kind string
}

func (e *PXEError) Error() string {
	return e.Msg
}

func newPXEError(kind, msg string, a ...interface{}) *PXEError {
	err := PXEError{
		Kind: kind,
		Msg:  fmt.Sprintf(msg, a...),
	}
	if kind == "TECHNICAL" {
		logger.Error(msg, err.Error())
	} else {
		logger.Info(msg, err.Error())
	}
	return &err
}

func ReadConfigurationContent(appConfig *model.AppConfig, name string) (*model.ConfigurationDetails, error) {
	logger.Info("Show configuration :: %s", name)

	var configToShow *model.Configuration
	for _, c := range ReadConfigurations(appConfig) {
		if name == c.Name {
			configToShow = c
			break
		}
	}
	if configToShow == nil {
		return nil, newPXEError("NOT_FOUND", "Configuration '%s' does not exists", name)
	}

	file := fmt.Sprintf("%s/%s/%s", appConfig.Configuration.Directory, configToShow.Bootloader.Name, configToShow.Name)

	f, err := ioutil.ReadFile(file)
	if err != nil {
		return nil, err
	}

	c := &model.ConfigurationDetails{
		Name:    name,
		Content: string(f),
	}

	return c, nil
}

func DeployConfiguration(appConfig *model.AppConfig, name string, hosts []*model.HostQuery) (*model.HostsResponse, error) {
	logger.Info("Deploy configuration :: %s :: %+v", name, hosts)

	var configToDeploy *model.Configuration
	for _, c := range ReadConfigurations(appConfig) {
		if name == c.Name {
			configToDeploy = c
			break
		}
	}
	if configToDeploy == nil {
		return nil, newPXEError("NOT_FOUND", "Configuration '%s' does not exists", name)
	}

	// Build maps in oder to optimize further searches
	hostsByName := make(map[string]*model.Host)
	hostsByPrimaryMAC := make(map[string]*model.Host)
	hostsByMAC := make(map[string]*model.Host)

	for _, h := range ReadHosts(appConfig, false) {
		hostsByName[h.Name] = h
		hostsByPrimaryMAC[h.MACAddresses[0]] = h
		for _, mac := range h.MACAddresses {
			hostsByMAC[strings.ToLower(mac)] = h
		}
	}

	hostsToDeploy := make(map[string]*model.Host)

	// Iterate over `hosts`
	for _, qh := range hosts {
		qh.MACAddress = strings.ToLower(qh.MACAddress)
		logger.Info("Processing :: %+v", qh)

		if qh.Configuration != "" {
			return nil, newPXEError("CONFILCT", "Configuration attribute for a host in this context is not allowed")
		}

		if qh.Name != "" {
			if hostsByName[qh.Name] == nil {
				return nil, newPXEError("NOT_FOUND", "No host declared for name <%s>", qh.Name)
			}
			if qh.MACAddress != "" {
				if hostsByMAC[qh.MACAddress] != nil {
					host := hostsByMAC[qh.MACAddress]
					if host.Name != qh.Name {
						return nil, newPXEError("CONFLICT", "Host <%s> does not match MAC address <%s>", qh.Name, qh.MACAddress)
					}
					if hostsToDeploy[host.Name] != nil {
						return nil, newPXEError("CONFLICT", "Host <%s> appears several times in query", host.Name)
					}
					hostsToDeploy[host.Name] = host

					continue
				}
				return nil, newPXEError("CONFLICT", "MAC address <%s> does not match host <%s>", qh.MACAddress, qh.Name)
			}

			if hostsToDeploy[qh.Name] != nil {
				return nil, newPXEError("CONFLICT", "Host <%s> appears several times in query", qh.Name)
			}
			hostsToDeploy[qh.Name] = hostsByName[qh.Name]

			continue
		}

		if qh.MACAddress != "" {
			if hostsByMAC[qh.MACAddress] == nil {
				return nil, newPXEError("NOT_FOUND", "No host declared with MAC address <%s>", qh.MACAddress)
			}

			host := hostsByMAC[qh.MACAddress]
			if hostsToDeploy[host.Name] != nil {
				return nil, newPXEError("CONFLICT", "Host <%s> appears several times in query", host.Name)
			}
			hostsToDeploy[host.Name] = host

			continue
		}

		return nil, newPXEError("BAD_REQUEST", "Either Name or MACAddress must be provided for each Host")
	}

	logger.Info("Host to deploy with configuration <%s> :: %+v", name, hostsToDeploy)
	for _, h := range hostsToDeploy {

		for _, mac := range h.MACAddresses {
			// Delete old config
			pxeFilePath := utils.PXEFilePathFromMAC(appConfig, mac)
			if _, err := os.Lstat(pxeFilePath); err == nil {
				os.Remove(pxeFilePath)
			}
		}

		// Create new config
		pxeFilePath := fmt.Sprintf("%s/pxelinux.cfg/%s", appConfig.Tftp.Root, utils.PXEFilenameFromMAC(h.MACAddresses[0]))
		srcConfigPath := fmt.Sprintf("%s/%s/%s", appConfig.Configuration.Directory, configToDeploy.Bootloader.Name, configToDeploy.Name)

		logger.Info("Creating symlink %s -> %s", srcConfigPath, pxeFilePath)
		err := os.Symlink(srcConfigPath, pxeFilePath)
		if err != nil {
			return nil, newPXEError("TECHNICAL", "Unable to create symlink")
		}

		for i := 1; i < len(h.MACAddresses); i++ {
			err := os.Symlink(pxeFilePath, fmt.Sprintf("%s/pxelinux.cfg/%s", appConfig.Tftp.Root, utils.PXEFilenameFromMAC(h.MACAddresses[i])))
			if err != nil {
				return nil, newPXEError("TECHNICAL", "Unable to create symlink")
			}
		}
	}

	hostsResponse := make([]*model.HostResponse, 0)

	resp := &model.HostsResponse{
		Hosts: hostsResponse,
	}

	for _, h := range hosts {
		hostResponse := &model.HostResponse{
			Name:          h.Name,
			Configuration: h.Configuration,
		}
		if h.Reboot {
			err := RebootHost(appConfig, hostsByName[h.Name])
			if err != nil {
				hostResponse.Rebooted = "ERROR"
			} else {
				hostResponse.Rebooted = "YES"
			}
		} else {
			hostResponse.Rebooted = "NO"
		}
		resp.Hosts = append(resp.Hosts, hostResponse)
	}

	return resp, nil
}

func getHostManagementAdapter(appConfig *model.AppConfig, host *model.Host) (HostManagementAdapter, error) {

	if host.IPMI != nil && host.Management != nil {
		return nil, logger.Errorf("Host '%s' has several host management adpter", host.Name)
	}

	var adapter HostManagementAdapter

	if host.IPMI != nil {
		adapter = IpmiBmcAdapter{}
		return adapter, nil
	}

	if host.Management != nil {
		for _, rma := range appConfig.HostManagementAdapters {
			if rma.Name == host.Management.AdapterName {
				host.Management.Adapter = rma
				break
			}
		}
		if host.Management.Adapter == nil {
			return nil, logger.Errorf("No host management adapter found with name '%s'", host.Management.AdapterName)
		}
		adapter = GenericHostManagementAdapter{}
		return adapter, nil
	}

	return nil, nil
}
