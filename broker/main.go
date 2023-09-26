package main

import (
	"broker/client"
	"broker/server"
	"os"
)

var (
	PORT = os.Getenv("PORT")
)

func main() {
	server.Run(PORT, client.New())
}
