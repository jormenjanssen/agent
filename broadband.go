package main

import (
	"fmt"
	"net"
)

//ConnectBroadband function.
func ConnectBroadband(address string, port int) (broadband Broadband, err error) {
	fmt.Printf("Connecting to: %s", address)
	tcpAdress := fmt.Sprintf("%s:%v", address, port)
	con, err := net.Dial("tcp", tcpAdress)
	if err != nil {
		return Broadband{Connection: nil}, err
	}

	return Broadband{Connection: con}, nil
}

// Broadband type.
type Broadband struct {
	Connection net.Conn
}
