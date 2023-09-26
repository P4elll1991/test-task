package server

import (
	"context"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha512"
	"log"
	"net"

	pb "github.com/p4elll1991/proto-for-test-task/broker"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
)

type client interface {
	EventBus(destination string, data []byte) error
}

type brokerServer struct {
	client client
	pb.UnimplementedBrokerServer
}

func (s *brokerServer) SendChallenge(ctx context.Context, req *pb.ChallengeRequest) (*pb.ChallengeResponse, error) {
	// Генерация 32 случайных байт
	data := make([]byte, 32)
	if _, err := rand.Read(data); err != nil {
		log.Printf("sendChallenge: %s", err.Error())
		return nil, grpc.Errorf(codes.Unknown, "%s", err.Error())
	}

	if err := s.client.EventBus(req.Destination, data); err != nil {
		log.Printf("sendChallenge: %s", err.Error())
		return nil, grpc.Errorf(grpc.Code(err), "%s", err.Error())
	}

	return &pb.ChallengeResponse{}, nil
}

func (s *brokerServer) SendMessage(ctx context.Context, req *pb.MessageRequest) (*pb.MessageResponse, error) {

	// Генерация случайного RSA ключа
	privateKey, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		log.Printf("sendMessage: %s", err.Error())
		return nil, grpc.Errorf(codes.Unknown, "%s", err.Error())
	}

	// Шифрование данных с использованием открытого ключа
	publicKey := &privateKey.PublicKey
	encryptedData, err := rsa.EncryptOAEP(
		sha512.New(),
		rand.Reader,
		publicKey,
		req.Data,
		nil,
	)
	if err != nil {
		log.Printf("sendMessage: %s", err.Error())
		return nil, grpc.Errorf(codes.Unknown, "%s", err.Error())
	}

	if err := s.client.EventBus(req.Destination, encryptedData); err != nil {
		log.Printf("sendMessage: %s", err.Error())
		return nil, grpc.Errorf(grpc.Code(err), "%s", err.Error())
	}

	return &pb.MessageResponse{}, nil
}

func Run(port string, client client) {
	lis, err := net.Listen("tcp", port)
	if err != nil {
		log.Print(err)
		return
	}
	s := grpc.NewServer()
	pb.RegisterBrokerServer(s, &brokerServer{client: client})
	if err := s.Serve(lis); err != nil {
		log.Print(err)
		return
	}
}
