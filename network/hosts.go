package network

import (
	"fmt"
	"io/ioutil"
	"net"
	"strconv"
	"strings"
)

type Host struct {
	address net.IP
	port    uint32
}

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
		if len(strings.Trim(host, " ")) == 0 {
			continue
		}

		addressStr, portStr, err := net.SplitHostPort(host)
		if err != nil {
			return nil, err
		}

		address := net.ParseIP(addressStr)

		port, err := strconv.Atoi(portStr)
		if err != nil {
			return nil, err
		}

		cleanHosts = append(cleanHosts, Host{address, uint32(port)})
	}

	return cleanHosts, nil
}

// Adds a new host (address:port) to the list of known hosts
func AddHost(address net.IP, port uint32, completeAddress string) error {
	hosts, err := GetHosts(completeAddress)
	if err != nil {
		return err
	}

	newHost := Host{address, port}
	for _, h := range hosts {
		// If host already exists, do nothing
		// Could return error, but seems unnecessary given that it just means that it shouldn't
		// be added again.
		if h.address.String() == newHost.address.String() && h.port == newHost.port {
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
		data += fmt.Sprintf("%s:%d\n", host.address.String(), host.port)
	}

	return ioutil.WriteFile(path, []byte(data), 0644)
}
