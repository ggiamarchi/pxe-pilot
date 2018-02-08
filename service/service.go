package service

import (
	"io/ioutil"
	"os"
	"strings"
	"sync"

	"fmt"

	"github.com/ggiamarchi/pxe-pilot/common/utils"
	"github.com/ggiamarchi/pxe-pilot/logger"
	"github.com/ggiamarchi/pxe-pilot/model"
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

	var wg sync.WaitGroup

	for i, host := range appConfig.Hosts {

		if host.IPMI != nil {
			wg.Add(1)
			hostlocal := host
			go func() {
				defer wg.Done()
				ChassisPowerStatus(hostlocal.IPMI)
			}()
		}

		hosts[i] = &model.Host{
			Name:         host.Name,
			MACAddresses: host.MACAddresses,
			IPMI:         host.IPMI,
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

	wg.Wait()

	return hosts
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
		return newPXEError("NOT_FOUND", "Configuration '%s' does not exists", name)
	}

	// Build maps in oder to optimize further searches
	hostsByName := make(map[string]*model.Host)
	hostsByPrimaryMAC := make(map[string]*model.Host)
	hostsByMAC := make(map[string]*model.Host)

	for _, h := range ReadHosts(appConfig) {
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
			return newPXEError("CONFILCT", "Configuration attribute for a host in this context is not allowed")
		}

		if qh.Name != "" {
			if hostsByName[qh.Name] == nil {
				return newPXEError("NOT_FOUND", "No host declared for name <%s>", qh.Name)
			}
			if qh.MACAddress != "" {
				if hostsByMAC[qh.MACAddress] != nil {
					host := hostsByMAC[qh.MACAddress]
					if host.Name != qh.Name {
						return newPXEError("CONFLICT", "Host <%s> does not match MAC address <%s>", qh.Name, qh.MACAddress)
					}
					if hostsToDeploy[host.Name] != nil {
						return newPXEError("CONFLICT", "Host <%s> appears several times in query", host.Name)
					}
					hostsToDeploy[host.Name] = host

					continue
				}
				return newPXEError("CONFLICT", "MAC address <%s> does not match host <%s>", qh.MACAddress, qh.Name)
			}

			if hostsToDeploy[qh.Name] != nil {
				return newPXEError("CONFLICT", "Host <%s> appears several times in query", qh.Name)
			}
			hostsToDeploy[qh.Name] = hostsByName[qh.Name]

			continue
		}

		if qh.MACAddress != "" {
			if hostsByMAC[qh.MACAddress] == nil {
				return newPXEError("NOT_FOUND", "No host declared with MAC address <%s>", qh.MACAddress)
			}

			host := hostsByMAC[qh.MACAddress]
			if hostsToDeploy[host.Name] != nil {
				return newPXEError("CONFLICT", "Host <%s> appears several times in query", host.Name)
			}
			hostsToDeploy[host.Name] = host

			continue
		}

		return newPXEError("BAD_REQUEST", "Either Name or MACAddress must be provided for each Host")
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
			return newPXEError("TECHNICAL", "Unable to create symlink")
		}

		for i := 1; i < len(h.MACAddresses); i++ {
			err := os.Symlink(pxeFilePath, fmt.Sprintf("%s/pxelinux.cfg/%s", appConfig.Tftp.Root, utils.PXEFilenameFromMAC(h.MACAddresses[i])))
			if err != nil {
				return newPXEError("TECHNICAL", "Unable to create symlink")
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
