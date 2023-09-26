package main

import (
	"http-server/client"
	"http-server/server"
	"os"
)

var (
	broker_addr = os.Getenv("BROKER_ADDR") // адрес grpc клиента, в который будут перенаправляться запросы
	PORT        = os.Getenv("PORT")        // порт на котором работает http сервер
)

func main() {
	// создаем объект grpc клиента
	client := client.New(broker_addr)
	// запускаем http сервер
	server.New(PORT, client).Run()
}
