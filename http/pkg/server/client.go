package server

import (
	"github.com/izveigor/p2p-words/http/pkg/config"
	"github.com/izveigor/p2p-words/http/pkg/server/pb"

	"google.golang.org/grpc"
)

type Client struct {
	Service pb.HTTPClient
}

func InitHTTPServiceClient() pb.HTTPClient {
	conn, err := grpc.Dial(config.Config.P2PSvcUrl, grpc.WithInsecure())
	if err != nil {
		panic(err)
	}

	return pb.NewHTTPClient(conn)
}

var HTTPServiceClient Client = Client{
	Service: InitHTTPServiceClient(),
}
