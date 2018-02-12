package main

import (
	"bufio"
	"fmt"
	"net"
	"strings"
	"time"
)

const Sep = '\r'

//ConnectBroadband function.
func ConnectBroadband(address string, port int) (broadband *Broadband, err error) {
	fmt.Printf("Connecting to: %s", address)
	tcpAdress := fmt.Sprintf("%s:%v", address, port)
	con, err := net.Dial("tcp", tcpAdress)
	if err != nil {
		return &Broadband{}, err
	}

	channel := make(chan AT, 1)
	broadband = &Broadband{at: channel}
	go handleInternal(con, broadband)

	return broadband, nil
}

// Broadband type.
type Broadband struct {
	at chan AT
}

// AT Command.
func (broadband Broadband) AT() *AT {

	at := AT{Command: "AT",
		ResponseChannel: make(chan string),
		HandledChannel:  make(chan string)}
	broadband.at <- at
	return &at
}

// Close the channel
func (broadband Broadband) Close() {
	close(broadband.at)
}

func (broadband Broadband) handle() {

}

func handleInternal(conn net.Conn, broadband *Broadband) {

	reader := bufio.NewReader(conn)
	timeoutDuration := 5 * time.Second

	for {

		message, more := <-broadband.at

		conn.SetWriteDeadline(time.Now().Add(timeoutDuration))

		_, err := conn.Write([]byte(message.Command + "\r"))

		if err != nil {
			fmt.Printf("failed to send: %v", err)
			close(message.ResponseChannel)
			continue

		}

		if !more {
			continue
		}

		exitCommand := false

		for {
			// Poll the channel if we're close/handled then stop and proces the next command.
			select {

			case _, err := <-message.HandledChannel:
				// Close the message
				fmt.Printf("closed channel: %v", err)
				close(message.ResponseChannel)
				exitCommand = true
			default:

				// Set max timeout
				conn.SetReadDeadline(time.Now().Add(timeoutDuration))
				// Read tokens delimited by newline
				str, err := reader.ReadString(Sep)
				if err != nil {
					fmt.Println(err)
					exitCommand = true
				}

				// If we have vallid data then redirect it to the caller.
				if !exitCommand {
					fmt.Printf("%s", str)

					// Redirect the string to the response stream.
					message.ResponseChannel <- strings.TrimSpace(str)
				}
			}

			// We need to exit here .. (avoid a panic with multiple close)
			if exitCommand {
				break
			}

		}
	}
}
