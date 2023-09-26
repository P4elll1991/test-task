package client

import (
	"context"
	"time"

	pb "github.com/p4elll1991/proto-for-test-task/broker"

	"google.golang.org/grpc"
)

type client struct {
	addr string
}

type connection struct {
	grpc   *grpc.ClientConn
	ctx    context.Context
	cancel context.CancelFunc
}

func newConnection(addr string) (*connection, error) {
	conn, err := grpc.Dial(addr, grpc.WithInsecure())
	if err != nil {
		return nil, err
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)

	return &connection{conn, ctx, cancel}, nil
}

func New(addr string) *client {
	return &client{addr: addr}
}

func (client *client) SendChallenge(destination string) error {
	conn, err := newConnection(client.addr)
	if err != nil {
		return err
	}
	defer conn.grpc.Close()
	defer conn.cancel()

	_, err = pb.NewBrokerClient(conn.grpc).SendChallenge(conn.ctx, &pb.ChallengeRequest{
		Destination: destination,
	})
	return err
}

func (client *client) SendMessage(destination string, data []byte) error {
	conn, err := newConnection(client.addr)
	if err != nil {
		return err
	}
	defer conn.grpc.Close()
	defer conn.cancel()

	_, err = pb.NewBrokerClient(conn.grpc).SendMessage(conn.ctx, &pb.MessageRequest{
		Destination: destination,
		Data:        data,
	})

	return err
}
