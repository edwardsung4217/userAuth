package main

import (
	"bufio"
	"fmt"
	"net"
	"os"
)

func main() {
	// Connect to the server
	conn, err := net.Dial("tcp", "localhost:3000") // Replace with your server's address
	if err != nil {
		fmt.Println("Failed to connect to the server:", err)
		return
	}
	defer conn.Close()

	reader := bufio.NewReader(os.Stdin)
	for {
		message, _ := reader.ReadString('\n')
		message = message[:len(message)-1]

		// Simulate sending data to the server
		_, err := conn.Write([]byte(message))

		if err != nil {
			fmt.Println("Failed to send data to the server:", err)
			return
		}
	}
}
