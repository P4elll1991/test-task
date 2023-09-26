package client

import (
	"context"
	"time"

	pb "github.com/p4elll1991/proto-for-test-task/broker"

	"google.golang.org/grpc"
)

// client - gRPC клиента
type client struct {
	addr string
}

// connection - объект соединения
type connection struct {
	grpc   *grpc.ClientConn
	ctx    context.Context
	cancel context.CancelFunc
}

// newConnection - устанавливает соединение по addr, а так же создает контекст
func newConnection(addr string) (*connection, error) {
	conn, err := grpc.Dial(addr, grpc.WithInsecure())
	if err != nil {
		return nil, err
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)

	return &connection{conn, ctx, cancel}, nil
}

func (conn *connection) Close() {
	conn.cancel()
	conn.grpc.Close()
}

func New(addr string) *client {
	return &client{addr: addr}
}

// SendChallenge - вызывает у gRPC клиента sendChallenge и передает туда destination
func (client *client) SendChallenge(destination string) error {
	conn, err := newConnection(client.addr)
	if err != nil {
		return err
	}
	defer conn.Close()

	_, err = pb.NewBrokerClient(conn.grpc).SendChallenge(conn.ctx, &pb.ChallengeRequest{
		Destination: destination,
	})
	return err
}

// SendMessage - вызывает у gRPC клиента sendMessage и передает туда destination и data
func (client *client) SendMessage(destination string, data []byte) error {
	conn, err := newConnection(client.addr)
	if err != nil {
		return err
	}
	defer conn.Close()

	_, err = pb.NewBrokerClient(conn.grpc).SendMessage(conn.ctx, &pb.MessageRequest{
		Destination: destination,
		Data:        data,
	})

	return err
}
