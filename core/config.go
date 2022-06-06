package core

import (
	"fmt"
	"net"
	"peercloud/crypto"
)

type Config struct {
	Address net.IP
	Port    uint16
}

func (c *Config) GetCompleteAddress() string {
	return fmt.Sprintf("%s:%d", c.Address, c.Port)
}

func (c *Config) GetNodeIdentifier() string {
	ip := c.GetCompleteAddress()

	identifier := crypto.HashSha256([]byte(ip))
	return crypto.HashAsString(identifier)
}
