package main

import (
	"os"
)

func main() {

	if os.Args == nil || os.Args[1] == "client" || os.Args[1] == "" {

		IP := "127.0.0.2"
		PORT := "9090"
		if len(os.Args) >= 4 {
			IP = os.Args[2]
			PORT = os.Args[3]
		}

		client := cInit()
		client.cStart(IP, PORT)
	} else {
		server := Init("tcp", "127.0.0.2", "9090")
		server.Start()
	}

}
