package main

import (
	"fmt"
	"userAuth/types"
)

func main() {
	var user types.User
	user.Username = "123"
	user.Password = "123"

	server := types.NewServer(":3000", user)

	// Access the msgch channel using the GetMsgChannel method
	msgch := server.GetMsgChannel()

	go func() {
		for payload := range msgch {
			// Access the user field using the GetUser method
			user := payload.GetAddress()
			fmt.Printf("Received from %s: %s\n", user, string(payload.GetMessage()))
		}
	}()

	if err := server.Start(); err != nil {
		fmt.Println("Server error:", err)
	}
}
