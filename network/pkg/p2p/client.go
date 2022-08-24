package p2p

import (
	"github.com/izveigor/p2p-words/network/pkg/config"
	"github.com/izveigor/p2p-words/network/pkg/p2p/pb"

	"google.golang.org/grpc"
)

type Client struct {
	Service pb.LemmatizersClient
}

func InitLemmatizersServiceClient() pb.LemmatizersClient {
	conn, err := grpc.Dial(config.Config.LemmatizerSvcUrl, grpc.WithInsecure())
	if err != nil {
		panic(err)
	}

	return pb.NewLemmatizersClient(conn)
}

var LemmatizersServiceClient Client = Client{
	Service: InitLemmatizersServiceClient(),
}
