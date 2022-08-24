package p2p

import (
	"context"
	"net"

	"github.com/izveigor/p2p-words/network/pkg/p2p/pb"

	"google.golang.org/grpc"
)

type server struct {
	pb.UnimplementedHTTPServer
}

func (s *server) SearchSentences(ctx context.Context, in *pb.SearchSentencesRequest) (*pb.SearchSentencesResponse, error) {
	var word string = in.GetWord()

	var sentences []string = append(
		GetInformation(word),
		UserLibrary.SearchSentences(word)...,
	)

	return &pb.SearchSentencesResponse{
		Sentences: sentences,
	}, nil
}

func (s *server) CreateBook(ctx context.Context, in *pb.CreateBookRequest) (*pb.CreateBookResponse, error) {
	CreateFiles(in.GetName(), in.GetText())
	return &pb.CreateBookResponse{Ok: true}, nil
}

func (s *server) GetBooksInformation(ctx context.Context, in *pb.GetBooksInformationRequest) (*pb.GetBooksInformationResponse, error) {
	return UserLibrary.GetBooksInformation(), nil
}

func InitHTTPServiceServer(address string) {
	lis, err := net.Listen("tcp", address)
	if err != nil {
		panic(err)
	}

	s := grpc.NewServer()
	pb.RegisterHTTPServer(s, &server{})
	if err := s.Serve(lis); err != nil {
		panic(err)
	}
}
