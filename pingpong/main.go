package main

import (
	"fmt"
	"strings"
)

func shout(ping chan string, pong chan string) {
	for {
		dataIn := <-ping

		pong <- fmt.Sprintf("%s!!!", strings.ToUpper(dataIn))
	}

}

func main() {
	fmt.Println("main")

	ping := make(chan string)
	pong := make(chan string)

	go shout(ping, pong)
	fmt.Println("Type something and press ENTER (enter Q to quit)")

	for {
		// print a prompt
		fmt.Print("-> ")

		// get user input
		var userInput string
		_, _ = fmt.Scanln(&userInput)

		if userInput == strings.ToLower("q") {
			// jump out of for loop
			break
		}

		// send userInput to "ping" channel
		ping <- userInput

		// wait for a response from the pong channel. Again, program
		// blocks (pauses) until it receives something from
		// that channel.
		response := <-pong

		// print the response to the console.
		fmt.Println("Response:", response)
	}

	fmt.Println("All done. Closing channels.")

	// close the channels
	close(ping)
	close(pong)
}
