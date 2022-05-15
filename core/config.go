package core

import (
	"fmt"
	"peercloud/crypto"
)

type Config struct {
	Address string
	Port    uint32
}

func (c *Config) GetCompleteAddress() string {
	return fmt.Sprintf("%s:%d", c.Address, c.Port)
}

func (c *Config) GetNodeIdentifier() string {
	ip := c.GetCompleteAddress()

	identifier := crypto.HashSha256([]byte(ip))
	return crypto.HashAsString(identifier)
}
