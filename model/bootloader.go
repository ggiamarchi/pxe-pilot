package model

import (
	"fmt"
)

type Bootloader struct {
	Name       string `json:"name" yaml:"name"`
	File       string `json:"file" yaml:"file"`
	ConfigPath string `json:"config_path" yaml:"config_path"`
}

func (b *Bootloader) String() string {
	return fmt.Sprintf("%+v", *b)
}
