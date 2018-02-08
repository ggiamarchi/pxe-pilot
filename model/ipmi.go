package model

import "fmt"

type IPMI struct {
	MACAddress string `json:"macAddress" yaml:"mac_address"`
	Username   string `json:"username" yaml:"username"`
	Password   string `json:"password" yaml:"password"`
	Interface  string `json:"interface" yaml:"interface"`
	Status     string `json:"status" yaml:"status"`
	Hostname   string `json:"hostname" yaml:"hostname"`
}

func (i *IPMI) String() string {
	return fmt.Sprintf("%+v", *i)
}
