package main

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"net"
	"strings"
	"time"
)

// SupplicantState message structure
type SupplicantState struct {
	LastInterval time.Time
	Available    bool
	Connected    bool
	BSSID        string
	KeyType      string
}

// Supplicant structure
type Supplicant struct {
	rwc io.ReadWriteCloser
}

// HandleAsync supplicant connection
func (supplicant *Supplicant) HandleAsync(supplicantChannel chan<- SupplicantState) {
	go supplicant.handleInternal(supplicantChannel)
}

//EncodeCommand converts string command to byte aray
func EncodeCommand(command string) []byte {

	cmdData := []byte(command)
	cmdData = cmdData[:len(cmdData)+1]

	copy(cmdData, []byte(command))

	return cmdData
}

func (supplicant *Supplicant) handleInternal(supplicantChannel chan<- SupplicantState) {

	bfr := bufio.NewReader(supplicant.rwc)

	// Close when we're done.
	defer supplicant.rwc.Close()

	timeout := 100 * time.Millisecond

	recvbuf := make([]byte, 4096)

	pingcmd := EncodeCommand("STATUS")

	for {

		supplicant.rwc.Write(pingcmd)

		n, err := bfr.Read(recvbuf)

		if err != nil {
			fmt.Println(err.Error())
			return
		}

		recvfinalbuf := make([]byte, n-1)
		copy(recvfinalbuf, recvbuf)

		str := (string(recvfinalbuf))

		fmt.Printf("received: %v\n", str)

		supplicantChannel <- SupplicantState{Available: true, LastInterval: time.Now()}

		time.Sleep(timeout)
	}

}

// Init function for the supplicant
func Init(address string) (supplicant *Supplicant, err error) {

	// Split up the protocol and the address eg. tcp:192.168.85.80 or unix:/var/run/wpa_supplicant/wifi0
	split := strings.Split(address, ":")
	innerAddress := fmt.Sprintf("%s:%s", split[1], split[2])

	// Connect to the supplicant.
	connection, err := ConnectSupplicant(split[0], innerAddress)

	if err != nil {
		return nil, err
	}

	// Return by pointer
	return &Supplicant{rwc: connection}, err
}

// ConnectSupplicant function
func ConnectSupplicant(protocol string, address string) (rwc io.ReadWriteCloser, err error) {

	// Simple abstraction divider
	switch protocol {

	case "tcp":
		{
			tcp, err := connectTCP(address)
			return tcp, err
		}
	default:
		{
			return nil, errors.New("nil")
		}
	}

}

func connectTCP(address string) (rwc io.ReadWriteCloser, err error) {
	return net.Dial("tcp", address)
}
