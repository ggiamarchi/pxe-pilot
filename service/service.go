package service

import (
	"io/ioutil"
	"os"
	"strings"

	"fmt"

	"dev.splitted-desktop.com/horizon/pxe-pilot/common/utils"
	"dev.splitted-desktop.com/horizon/pxe-pilot/logger"
	"dev.splitted-desktop.com/horizon/pxe-pilot/model"
)

func ReadConfigurations(appConfig *model.AppConfig) []*model.Configuration {
	files, _ := ioutil.ReadDir(appConfig.Configuration.Directory)
	configurations := make([]*model.Configuration, len(files))
	for i, f := range files {
		configurations[i] = &model.Configuration{
			Name: f.Name(),
		}
	}
	return configurations
}

func ReadHosts(appConfig *model.AppConfig) []*model.Host {

	pxelinuxDir := appConfig.Tftp.Root + "/pxelinux.cfg"

	files, _ := ioutil.ReadDir(pxelinuxDir)

	filenames := make([]string, len(files))
	for i, f := range files {
		filenames[i] = f.Name()
	}

	hosts := make([]*model.Host, len(appConfig.Hosts))

	for i, host := range appConfig.Hosts {
		hosts[i] = &model.Host{
			Name:         host.Name,
			MACAddresses: host.MACAddresses,
		}
		pxeFile := utils.PXEFilenameFromMAC(hosts[i].MACAddresses[0])
		pxeFilePath := fmt.Sprintf("%s/%s", pxelinuxDir, pxeFile)

		if _, err := os.Stat(pxeFilePath); err != nil {
			// No PXE config deployed for this host
			continue
		}

		configFile, err := os.Readlink(pxeFilePath)
		if err != nil {
			panic(err)
		}
		hosts[i].Configuration = &model.Configuration{
			Name: configFile[strings.LastIndex(configFile, "/")+1:],
		}
	}

	return hosts
}

func DeployConfiguration(appConfig *model.AppConfig, name string, hosts []*model.HostQuery) error {
	logger.Info("Deploy configuration :: %s :: %+v", name, hosts)

	configExists := false
	for _, c := range ReadConfigurations(appConfig) {
		if name == c.Name {
			configExists = true
			break
		}
	}
	if !configExists {
		return logger.Errorf("Configuration '%s' does not exists", name)
	}
	//	configFile := fmt.Sprintf("%s/%s", appConfig.Configuration.Directory, name)

	// Build maps in oder to optimize further searches
	hostsByName := make(map[string]*model.Host)
	hostsByPrimaryMAC := make(map[string]*model.Host)
	hostsByMAC := make(map[string]*model.Host)

	for _, h := range ReadHosts(appConfig) {
		hostsByName[h.Name] = h
		hostsByPrimaryMAC[h.MACAddresses[0]] = h
		for _, mac := range h.MACAddresses {
			hostsByMAC[mac] = h
		}
	}

	hostsToDeploy := make(map[string]*model.Host)

	// 2. Iterate over `hosts`
	for _, qh := range hosts {
		logger.Info("Processing :: %+v", qh)

		if qh.Configuration != "" {
			return logger.Errorf("Configuration attribute for a host in this context is not allowed")
		}

		if qh.Name != "" {
			if hostsByName[qh.Name] == nil {
				return logger.Errorf("No host declared for name <%s>", qh.Name)
			}
			if qh.MACAddress != "" {
				if hostsByMAC[qh.MACAddress] != nil {
					host := hostsByMAC[qh.MACAddress]
					if host.Name != qh.Name {
						return logger.Errorf("Host <%s> does not match MAC address <%s>", qh.Name, qh.MACAddress)
					}
					if hostsToDeploy[host.Name] != nil {
						return logger.Errorf("Host <%s> appears several times in query", host.Name)
					}
					hostsToDeploy[host.Name] = host

					continue
				}
				return logger.Errorf("MAC address <%s> does not match host <%s>", qh.MACAddress, qh.Name)
			}

			if hostsToDeploy[qh.Name] != nil {
				return logger.Errorf("Host <%s> appears several times in query", qh.Name)
			}
			hostsToDeploy[qh.Name] = hostsByName[qh.Name]

			continue
		}

		if qh.MACAddress != "" {
			if hostsByMAC[qh.MACAddress] == nil {
				return logger.Errorf("No host declared with MAC address <%s>", qh.MACAddress)
			}

			host := hostsByMAC[qh.MACAddress]
			if hostsToDeploy[host.Name] != nil {
				return logger.Errorf("Host <%s> appears several times in query", host.Name)
			}
			hostsToDeploy[host.Name] = host

			continue
		}

		return logger.Errorf("Either Name or MACAddress must be provided for each Host")
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
		srcConfigPath := fmt.Sprintf("%s/%s", appConfig.Configuration.Directory, name)

		logger.Info("Creating symlink %s -> %s", srcConfigPath, pxeFilePath)
		err := os.Symlink(srcConfigPath, pxeFilePath)
		if err != nil {
			return logger.Errorf("Unable to create symlink")
		}

		for i := 1; i < len(h.MACAddresses); i++ {
			err := os.Symlink(pxeFilePath, fmt.Sprintf("%s/pxelinux.cfg/%s", appConfig.Tftp.Root, utils.PXEFilenameFromMAC(h.MACAddresses[i])))
			if err != nil {
				return logger.Errorf("Unable to create symlink")
			}
		}
	}

	return nil
}

func configDeployedForHost(a string, host *model.Host) bool {
	for _, b := range host.MACAddresses {
		if b == a {
			return true
		}
	}
	return false
}
