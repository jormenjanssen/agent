package main

// AT command structure
type AT struct {
	Command         string
	ResponseChannel chan string
	HandledChannel  chan string
}

// Handled closes the channel
func (at AT) Handled() {
	close(at.HandledChannel)
}
