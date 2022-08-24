package p2p

import (
	"github.com/izveigor/p2p-words/network/pkg/config"
	"github.com/izveigor/p2p-words/network/pkg/p2p/pb"
	"github.com/stretchr/testify/assert"

	"google.golang.org/grpc"

	"context"
	"net"
	"strconv"
	"testing"
	"time"
)

func TestPeer(t *testing.T) {
	peers = PeersInformation{}
	defer func() {
		peers = PeersInformation{}
	}()
	addresses := map[[5]Peer][]Peer{
		{
			Peer{Address: "a", Client: nil},
			Peer{Address: "c", Client: nil},
			Peer{Address: "d", Client: nil},
			Peer{Address: "b", Client: nil},
			Peer{Address: "e", Client: nil},
		}: {
			Peer{Address: "a", Client: nil},
			Peer{Address: "b", Client: nil},
			Peer{Address: "c", Client: nil},
			Peer{Address: "d", Client: nil},
			Peer{Address: "e", Client: nil},
		},
		{
			Peer{Address: "c", Client: nil},
			Peer{Address: "a", Client: nil},
			Peer{Address: "b", Client: nil},
			Peer{Address: "d", Client: nil},
			Peer{Address: "e", Client: nil},
		}: {
			Peer{Address: "a", Client: nil},
			Peer{Address: "b", Client: nil},
			Peer{Address: "c", Client: nil},
			Peer{Address: "d", Client: nil},
			Peer{Address: "e", Client: nil},
		},
	}
	for input, answer := range addresses {
		for _, peer := range input {
			peers.Add(peer)
		}
		assert.Equal(t, peers.Peers, answer)

		for i := 0; i < len(answer); i++ {
			peers.Delete(answer[i].Address)
			assert.Equal(t, peers.Peers, answer[i+1:])
		}
		peers = PeersInformation{}
	}
}

func startServer(t *testing.T, address string) {
	lis, err := net.Listen("tcp", address)
	if err != nil {
		t.Fatal(err)
	}

	grpcServer := grpc.NewServer()
	pb.RegisterP2PServer(grpcServer, &serverP2P{})
	go grpcServer.Serve(lis)
	time.Sleep(1250)
}

func getClient(t *testing.T, address string) pb.P2PClient {
	conn, err := grpc.Dial(address, grpc.WithInsecure())
	if err != nil {
		t.Fatal(err)
	}

	return pb.NewP2PClient(conn)
}

func TestConnect(t *testing.T) {
	myAddress = "localhost:3334"
	defer func() {
		myAddress = ""
	}()
	var address string = "localhost:3333"
	startServer(t, address)
	client := getClient(t, address)
	response, err := client.Connect(context.Background(), &pb.ConnectRequest{
		Address: address,
	})

	if err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, response.Address, myAddress)
}

func TestList(t *testing.T) {
	peers = PeersInformation{}
	peers.Add(Peer{Address: "localhost:3333", Client: nil})
	peers.Add(Peer{Address: "localhost:3335", Client: nil})
	defer func() {
		peers = PeersInformation{}
	}()
	var address string = "localhost:3334"
	startServer(t, address)
	client := getClient(t, address)
	response, err := client.List(context.Background(), &pb.ListRequest{})

	if err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, response.Addresses, []string{
		"localhost:3333",
		"localhost:3335",
	})
}

func TestDisconnect(t *testing.T) {
	peers = PeersInformation{}
	peers.Add(Peer{Address: "localhost:3333", Client: nil})
	peers.Add(Peer{Address: "localhost:3335", Client: nil})
	peers.Add(Peer{Address: "localhost:3334", Client: nil})
	defer func() {
		peers = PeersInformation{}
	}()

	var address string = "localhost:3335"
	startServer(t, address)
	client := getClient(t, address)
	_, err := client.Disconnect(context.Background(), &pb.DisconnectRequest{
		Address: address,
	})

	if err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, peers.Peers, []Peer{
		Peer{Address: "localhost:3333", Client: nil},
		Peer{Address: "localhost:3334", Client: nil},
	})
}

func TestGet(t *testing.T) {
	firstBook := &Book{
		Name:            "Название",
		CharactersCount: 200,
		WordsCount:      100,
		Tree: &Tree{
			Root: NewNode(&TreeElement{
				Word:      "слово",
				Sentences: []string{"Предложение"},
			}),
		},
	}
	UserLibrary = Library{[]*Book{}}
	UserLibrary.Books = append(UserLibrary.Books, firstBook)
	defer func() {
		UserLibrary = Library{[]*Book{}}
	}()

	var address string = "localhost:3336"
	startServer(t, address)
	client := getClient(t, address)
	response, err := client.Get(context.Background(), &pb.GetRequest{
		Address: address,
		Word:    "Слово",
	})

	if err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, response.Sentences, []string{"Предложение"})
}

func TestCreateServer(t *testing.T) {
	myAddress = "address"
	defer func() {
		myAddress = ""
	}()
	go createServer(2601)
	time.Sleep(time.Millisecond * 1250)

	var address string = config.Config.P2PAddress + ":2601"
	client := getClient(t, address)
	response, err := client.Connect(context.Background(), &pb.ConnectRequest{
		Address: address,
	})

	if err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, response.Address, myAddress)
}

func TestCreateClients(t *testing.T) {
	peers = PeersInformation{}
	defer func() {
		peers = PeersInformation{}
	}()

	go createServer(2602)
	time.Sleep(time.Millisecond * 1250)
	go createClient(config.Config.P2PAddress + ":2602")
	time.Sleep(time.Millisecond * 1250)

	assert.Equal(t, peers.Peers[0].Address, config.Config.P2PAddress+":2602")
}

func TestGetInformation(t *testing.T) {
	firstBook := &Book{
		Name:            "Название",
		CharactersCount: 200,
		WordsCount:      100,
		Tree: &Tree{
			Root: NewNode(&TreeElement{
				Word:      "слово",
				Sentences: []string{"Предложение"},
			}),
		},
	}
	UserLibrary = Library{[]*Book{}}
	UserLibrary.Books = append(UserLibrary.Books, firstBook)
	defer func() {
		UserLibrary = Library{[]*Book{}}
	}()

	go createServer(2500)
	time.Sleep(time.Millisecond * 1250)
	go createClient(config.Config.P2PAddress + ":2500")
	time.Sleep(time.Millisecond * 2000)

	sentences := GetInformation("слово")
	assert.Equal(t, sentences, []string{"Предложение"})
}

func TestStartServer(t *testing.T) {
	peers = PeersInformation{}
	defer func() {
		peers = PeersInformation{}
	}()
	go StartServer()
	time.Sleep(time.Millisecond * 1250)
	go StartServer()
	time.Sleep(time.Millisecond * 1250)
	addresses := []string{}
	for _, peer := range peers.Peers {
		addresses = append(addresses, peer.Address)
	}

	var port string = strconv.Itoa(config.Config.InitialPort)
	assert.Equal(t, addresses, []string{
		"localhost:" + port,
		"localhost:" + port,
		"localhost:2001",
	})
}

func TestStopServer(t *testing.T) {
	go createServer(4500)
	time.Sleep(time.Millisecond * 1250)
	go createClient("localhost:4500")
	time.Sleep(time.Millisecond * 1250)
	StopServer()
}
