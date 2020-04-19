package model

import (
	"fmt"
)

type HostManagementInfo struct {
	Name        string               `yaml:"name"`
	Vars        []*HostManagementVar `json:"vars" yaml:"vars"`
	PowerStatus *HostManagementCmd   `yaml:"power_status"`
	PowerOn     *HostManagementCmd   `yaml:"power_on"`
	PowerOff    *HostManagementCmd   `yaml:"power_off"`
	PowerReset  *HostManagementCmd   `yaml:"power_reset"`
}

func (b *HostManagementInfo) String() string {
	return fmt.Sprintf("%+v", *b)
}

type HostManagementCmd struct {
	Cmd        string `yaml:"cmd"`
	PatternOn  string `yaml:"pattern_on"`
	PatternOff string `yaml:"pattern_off"`
}

func (b *HostManagementCmd) String() string {
	return fmt.Sprintf("%+v", *b)
}

type HostManagement struct {
	MacAddress  string               `json:"macAddress" yaml:"mac_address"`
	IPAddress   string               `json:"ipAddress" yaml:"ip_address"`
	AdapterName string               `json:"adapter" yaml:"adapter"`
	Adapter     *HostManagementInfo  `json:"-" yaml:"-"`
	Subnet      string               `json:"subnet" yaml:"subnet"`
	Vars        []*HostManagementVar `json:"vars" yaml:"vars"`
}

func (b *HostManagement) String() string {
	return fmt.Sprintf("%+v", *b)
}

type HostManagementVar struct {
	Name  string `json:"name"  yaml:"name"`
	Value string `json:"value" yaml:"value"`
}

func (b *HostManagementVar) String() string {
	return fmt.Sprintf("%+v", *b)
}
