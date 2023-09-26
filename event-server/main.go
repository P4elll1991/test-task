package main

import (
	"context"
	"net"
	"os"

	"google.golang.org/grpc"

	"log"

	pb "github.com/p4elll1991/proto-for-test-task/event-server"
)

// eventServer - gRPC сервер, обрабатывающий события
type eventServer struct {
	pb.UnimplementedEventServerServer
}

// EventBus - получает событие Event и логирует содержащуюся там data
func (s *eventServer) EventBus(ctx context.Context, req *pb.EventBusRequest) (*pb.EventBusResponse, error) {
	log.Printf("an event was received with data: %v\n", req.Event.Data)
	return &pb.EventBusResponse{}, nil
}

var PORT = os.Getenv("PORT")

func main() {
	lis, err := net.Listen("tcp", PORT)
	if err != nil {
		log.Fatal(err)
		return
	}
	s := grpc.NewServer()
	pb.RegisterEventServerServer(s, &eventServer{})
	if err := s.Serve(lis); err != nil {
		log.Fatal(err)
		return
	}
}
