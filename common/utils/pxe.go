package utils

import (
	"fmt"
	"strings"

	"github.com/ggiamarchi/pxe-pilot/model"
)

func PXEFilenameFromMAC(mac string) string {
	return strings.ToLower(fmt.Sprintf("01-%s-%s-%s-%s-%s-%s", mac[0:2], mac[3:5], mac[6:8], mac[9:11], mac[12:14], mac[15:17]))
}

func PXEFilePathFromMAC(appConfig *model.AppConfig, mac string) string {
	return fmt.Sprintf("%s/pxelinux.cfg/%s", appConfig.Tftp.Root, PXEFilenameFromMAC(mac))
}
