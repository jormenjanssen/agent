package main

import "fmt"
import "time"
import "runtime"

func main() {
	fmt.Println("Agent starting")

	supplicant, err := Init("tcp:127.0.0.1:7791")
	if err != nil {
		fmt.Printf("Error failed to connect: %v", err.Error())
		return
	}

	supplicantChannel := make(chan SupplicantState, 10)
	supplicant.HandleAsync(supplicantChannel)

	for {
		// Poll the descriptors, sleep if no data is available
		select {

		case state, more := <-supplicantChannel:

			if !more {
				break
			}
			fmt.Printf("Available state: %v @ %v \n", state.Available, state.LastInterval)

		default:
			time.Sleep(250 * time.Millisecond)
			runtime.GC()
		}
	}

}
