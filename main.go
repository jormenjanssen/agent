package main

import "fmt"

func main() {
	fmt.Println("Agent starting")
	broadband, err := ConnectBroadband("127.0.0.1", 9001)

	if err != nil {
		fmt.Printf("Could not connect %s\n", err)
		return
	}

	fmt.Printf("Connected to: %v", broadband)

}
