package main

import (
	"http-server/client"
	"http-server/server"
	"os"
)

var (
	broker_addr = os.Getenv("BROKER_ADDR")
	PORT        = os.Getenv("PORT")
)

func main() {
	client := client.New(broker_addr)
	server.New(PORT, client).Run()
}
