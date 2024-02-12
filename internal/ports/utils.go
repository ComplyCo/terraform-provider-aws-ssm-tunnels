package ports

import (
	"fmt"
	"math/rand"
	"net"
)

func FindOpenPort(lowerPort, upperPort int) (int, error) {
	if lowerPort < 0 || upperPort < 0 {
		return 0, fmt.Errorf("port range must be positive")
	}

	if lowerPort > upperPort {
		return 0, fmt.Errorf("lower port must be less than upper port")
	}

	if lowerPort > 65535 || upperPort > 65535 {
		return 0, fmt.Errorf("port range must be less than 65536")
	}

	for port := rand.Intn(upperPort-lowerPort+1) + lowerPort; port <= upperPort; port++ {
		address := fmt.Sprintf(":%d", port)
		listener, err := net.Listen("tcp", address)
		if err == nil {
			listener.Close()
			return port, nil // Port is open
		}
	}
	return 0, fmt.Errorf("no open port found in the range %d-%d", lowerPort, upperPort)
}
