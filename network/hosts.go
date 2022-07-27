package network

import (
	"fmt"
	"io/ioutil"
	"net"
	"peercloud/crypto"
	"strconv"
	"strings"
)

// ===== Should mode to models.go? ======
type Host struct {
	Address net.IP
	Port    uint16
}

func (h *Host) GetCompleteAddress() string {
	return fmt.Sprintf("%s:%d", h.Address, h.Port)
}

func (h *Host) GetNodeIdentifier() string {
	ip := h.GetCompleteAddress()

	identifier := crypto.HashSha256([]byte(ip))
	return crypto.HashAsString(identifier)
}

// ==============================

// TEMPORAL: Should use constants

func getHostsPath(completeAddress string) string {
	return ".peercloud/" + completeAddress + "/hosts"
}

////  ///// ////

/*
const DEFAULT_HOSTS_PATH = ".peercloud/hosts"
*/

// Returns all hosts known by the node
func GetHosts(completeAddress string) ([]Host, error) {
	path := getHostsPath(completeAddress)

	bytes, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}

	content := string(bytes)
	hosts := strings.Split(content, "\n")

	var cleanHosts []Host
	// Continue if string is empty
	for _, host := range hosts {
		trimmed := strings.Trim(host, " ")
		if len(trimmed) == 0 {
			continue
		}

		addressStr, portStr, err := net.SplitHostPort(trimmed)
		if err != nil {
			return nil, err
		}

		address := net.ParseIP(addressStr)

		port, err := strconv.Atoi(portStr)
		if err != nil {
			return nil, err
		}

		cleanHosts = append(cleanHosts, Host{address, uint16(port)})
	}

	return cleanHosts, nil
}

// Adds a new host (address:port) to the list of known hosts
func AddHost(address net.IP, port uint16, completeAddress string) error {
	hosts, err := GetHosts(completeAddress)
	if err != nil {
		return err
	}

	newHost := Host{address, port}
	for _, h := range hosts {
		// If host already exists, do nothing
		// Could return error, but seems unnecessary given that it just means that it shouldn't
		// be added again.
		if h.Address.String() == newHost.Address.String() && h.Port == newHost.Port {
			return nil
		}
	}

	hosts = append(hosts, newHost)
	return saveHosts(hosts, completeAddress)
}

func saveHosts(hosts []Host, completeAddress string) error {
	path := getHostsPath(completeAddress)

	var data string
	for _, host := range hosts {
		data += fmt.Sprintf("%s:%d\n", host.Address.String(), host.Port)
	}

	return ioutil.WriteFile(path, []byte(data), 0644)
}
