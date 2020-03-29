package model

import (
	"fmt"
)

type Configuration struct {
	Name string `json:"name"`
}

func (c *Configuration) String() string {
	return fmt.Sprintf("%+v", *c)
}

type ConfigurationDetails struct {
	Name    string `json:"name"`
	Content string `json:"content"`
}

func (c *ConfigurationDetails) String() string {
	return fmt.Sprintf("%+v", *c)
}
