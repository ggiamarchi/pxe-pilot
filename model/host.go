package model

import "fmt"

type Host struct {
	Name          string         `json:"name" yaml:"name"`
	MACAddresses  []string       `json:"macAddresses" yaml:"mac_addresses"`
	Configuration *Configuration `json:"configuration" yaml:"configuration"`
	IPMI          *IPMI          `json:"ipmi" yaml:"ipmi"`
}

func (h *Host) String() string {
	return fmt.Sprintf("%+v", *h)
}

type HostQuery struct {
	Name          string `json:"name"`
	MACAddress    string `json:"macAddress"`
	Configuration string `json:"configuration"`
}

func (h *HostQuery) String() string {
	return fmt.Sprintf("%+v", *h)
}

type HostsQuery struct {
	Hosts []*HostQuery `json:"hosts"`
}

func (h *HostsQuery) String() string {
	return fmt.Sprintf("%+v", *h)
}
