package platformsh

import (
	"fmt"
	"net"
	"os"
)

func NewListener() (net.Listener, error) {
	if socket, ok := os.LookupEnv("SOCKET"); ok {
		return net.Listen("unix", socket)
	}

	if port, ok := os.LookupEnv("PORT"); ok {
		addr := fmt.Sprintf("127.0.0.1:%s", port)
		return net.Listen("tcp", addr)
	}

	return nil, missingEnvironment("SOCKET", "PORT")
}