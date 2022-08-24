package p2p

import (
	"github.com/izveigor/p2p-words/network/pkg/config"
	"github.com/izveigor/p2p-words/network/pkg/p2p/pb"
	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc"

	"context"
	"errors"
	"net"
	"testing"
	"time"
)

type Server struct {
	pb.UnimplementedLemmatizersServer
}

func (s *Server) Lemmatize(ctx context.Context, in *pb.LemmatizerRequest) (*pb.LemmatizerResponse, error) {
	words := []*pb.LemmatizedWord{
		{Word: "Слово", Sentence: "Предложение", Id: 1},
	}
	if in.GetText() == "Привет!" {
		return &pb.LemmatizerResponse{Words: words}, nil
	}
	return nil, errors.New("Not Nil")
}

func TestLemmatize(t *testing.T) {
	lis, err := net.Listen("tcp", config.Config.LemmatizerSvcUrl)
	if err != nil {
		t.Fatal(err)
	}

	grpcServer := grpc.NewServer()
	pb.RegisterLemmatizersServer(grpcServer, &Server{})
	go grpcServer.Serve(lis)
	time.Sleep(time.Millisecond * 1250)
	response, err := LemmatizersServiceClient.Service.Lemmatize(context.Background(), &pb.LemmatizerRequest{Text: "Привет!"})
	if err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, response.Words[0].Word, "Слово")
	assert.Equal(t, response.Words[0].Sentence, "Предложение")
	assert.Equal(t, response.Words[0].Id, int32(1))
}
