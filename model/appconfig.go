package model

import (
	"fmt"
)

type AppConfig struct {
	Hosts []*Host
	Tftp  struct {
		Root string
	}
	Configuration struct {
		Directory   string
		Bootloaders []*Bootloader
	}
	Server struct {
		Port int
	}
	HostManagementAdapters []*HostManagementInfo `yaml:"host_management_adapters"`
}

func (c *AppConfig) String() string {
	return fmt.Sprintf("%+v", *c)
}
