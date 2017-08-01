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
