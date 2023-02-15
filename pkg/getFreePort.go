package pkg

import (
	"net"
	"strconv"
)

func GetFreePort() (string, error) {
	addr, err := net.ResolveTCPAddr("tcp", "localhost:0")
	if err != nil {
		return "0", err
	}

	l, err := net.ListenTCP("tcp", addr)
	if err != nil {
		return "0", err
	}

	defer l.Close()
	port := strconv.Itoa(l.Addr().(*net.TCPAddr).Port)

	return port, nil
}
