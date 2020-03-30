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
	Reboot        bool   `json:"reboot"`
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

type HostResponse struct {
	Name          string `json:"name"`
	Configuration string `json:"configuration"`
	Rebooted      string `json:"rebooted"`
}

func (h *HostResponse) String() string {
	return fmt.Sprintf("%+v", *h)
}

type HostsResponse struct {
	Hosts []*HostResponse `json:"hosts"`
}

func (h *HostsResponse) String() string {
	return fmt.Sprintf("%+v", *h)
}
