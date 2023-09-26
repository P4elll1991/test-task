package client

import (
	"context"
	"fmt"
	"time"

	pb "github.com/p4elll1991/proto-for-test-task/event-server"

	"google.golang.org/grpc"
)

type client struct {
	conn *grpc.ClientConn
}

type connection struct {
	grpc   *grpc.ClientConn
	ctx    context.Context
	cancel context.CancelFunc
}

func newConnection(addr string) (*connection, error) {
	conn, err := grpc.Dial(addr, grpc.WithInsecure())
	if err != nil {
		return nil, fmt.Errorf("Failed to connect: %v", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)

	return &connection{conn, ctx, cancel}, nil
}

func New() *client {
	return &client{}
}

func (client *client) EventBus(destination string, data []byte) error {
	conn, err := newConnection(destination)
	if err != nil {
		return err
	}
	defer conn.grpc.Close()
	defer conn.cancel()

	_, err = pb.NewEventServerClient(conn.grpc).EventBus(conn.ctx, &pb.EventBusRequest{
		Event: &pb.Event{
			Data: data,
		},
	})

	return err
}
