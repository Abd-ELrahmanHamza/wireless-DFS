package main

import (
	"fmt"
	"net"
)

func findThreeAvailablePorts(startPort, endPort int) ([]int, error) {
	var availablePorts []int

	for port := startPort; port <= endPort; port++ {
		addr := fmt.Sprintf("localhost:%d", port)
		listener, err := net.Listen("tcp", addr)
		if err != nil {
			// Port is not available, continue to the next port
			continue
		}
		// Port is available, close the listener and add it to the list of available ports
		listener.Close()
		availablePorts = append(availablePorts, port)
		// Stop iterating if three available ports are found
		if len(availablePorts) == 3 {
			break
		}
	}

	if len(availablePorts) == 0 {
		return nil, fmt.Errorf("no available ports found in the range %d-%d", startPort, endPort)
	}

	return availablePorts, nil
}
