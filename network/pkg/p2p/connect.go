package p2p

import (
	"context"
	"log"
	"net"
	"os"
	"strconv"
	"strings"
	"sync"

	"github.com/izveigor/p2p-words/network/pkg/config"
	"github.com/izveigor/p2p-words/network/pkg/p2p/pb"
	"google.golang.org/grpc"
	"google.golang.org/grpc/status"
)

type Peer struct {
	Address string
	Client  pb.P2PClient
}

type PeersInformation struct {
	Peers []Peer
}

func (p *PeersInformation) Add(peer Peer) {
	var index int = 0
	if len(p.Peers) == 0 {
		p.Peers = append(p.Peers, peer)
	} else {

		for p.Peers[index].Address < peer.Address && index < len(p.Peers)-1 {
			index++
		}

		if index == 0 && peer.Address < p.Peers[index].Address {
			p.Peers = append([]Peer{peer}, p.Peers...)
		} else if index == len(p.Peers)-1 && p.Peers[index].Address < peer.Address {
			p.Peers = append(p.Peers, peer)
		} else {
			p.Peers = append(p.Peers[:index+1], p.Peers[index:]...)
			p.Peers[index] = peer
		}
	}
}

func (p *PeersInformation) Delete(address string) {
	var index int = p.binarySearch(address)
	if index != -1 {
		p.Peers = append(p.Peers[:index], p.Peers[index+1:]...)
	}
}

func (p *PeersInformation) binarySearch(address string) int {
	var left, right int = 0, len(p.Peers) - 1

	for left < right {
		var mid int = (left + right) / 2
		if address == p.Peers[mid].Address {
			return mid
		} else if p.Peers[mid].Address < address {
			left = mid + 1
		} else {
			right = mid
		}
	}
	return left
}

var (
	peers            = PeersInformation{}
	maxPeers         = 2000
	myAddress string = ""
)

type serverP2P struct {
	pb.UnimplementedP2PServer
}

func (s *serverP2P) Connect(ctx context.Context, in *pb.ConnectRequest) (*pb.ConnectResponse, error) {
	var address string = in.GetAddress()
	if !strings.HasSuffix(os.Args[0], ".test") {
		if peers.Peers[len(peers.Peers)-1].Address != address {
			go createClient(address)
		} else {
			log.Printf("CONNECT: Пир с адресом %v присоединился к сети.", address)
		}
	} else {
		log.Printf("CONNECT: Пир с адресом %v присоединился к сети.", address)
	}
	return &pb.ConnectResponse{
		Address: myAddress,
	}, nil
}

func (s *serverP2P) List(ctx context.Context, in *pb.ListRequest) (*pb.ListResponse, error) {
	log.Printf("LIST: Отправил список пиров неизвестному пиру.")
	if len(peers.Peers) > maxPeers {
		return nil, status.Error(403, "Достигнуто максимальное количество пиров.")
	}
	addresses := []string{}
	for _, peer := range peers.Peers {
		addresses = append(addresses, peer.Address)
	}
	return &pb.ListResponse{
		Addresses: addresses,
	}, nil
}

func (s *serverP2P) Disconnect(ctx context.Context, in *pb.DisconnectRequest) (*pb.DisconnectResponse, error) {
	var address string = in.GetAddress()
	peers.Delete(address)
	log.Printf("DISCONNECT: Пир с адресом %v отсоединился от сети.", address)
	return &pb.DisconnectResponse{
		Address: myAddress,
	}, nil
}

func (s *serverP2P) Get(ctx context.Context, in *pb.GetRequest) (*pb.GetResponse, error) {
	var address, word string = in.GetAddress(), strings.ToLower(in.GetWord())
	log.Printf("GET: Пир с адресом %v ищет предложения со словом %q.", address, word)
	sentences := UserLibrary.SearchSentences(word)
	return &pb.GetResponse{
		Sentences: sentences,
	}, nil
}

func createServer(port int) {
	myAddress = config.Config.P2PAddress + ":" + strconv.Itoa(port)
	lis, err := net.Listen("tcp", myAddress)
	if err != nil {
		panic(err)
	}

	s := grpc.NewServer()
	pb.RegisterP2PServer(s, &serverP2P{})
	log.Printf("INFO: Начинаем запускать сервер по адресу %v", myAddress)
	peers.Add(Peer{
		Address: myAddress,
		Client:  nil,
	})
	if err := s.Serve(lis); err != nil {
		panic(err)
	}
}

func createClient(address string) {
	conn, err := grpc.Dial(address, grpc.WithInsecure())
	if err != nil {
		log.Printf("ERROR: Не смогли подключиться к адресу %v.", address)
	}

	client := pb.NewP2PClient(conn)
	response, err := client.Connect(context.Background(), &pb.ConnectRequest{
		Address: myAddress,
	})

	if err == nil && response.Address == address {
		log.Printf("CONNECT: Установлено соединение с пиром %v", response.Address)
	} else {
		log.Printf("ERROR: Не смогли подключиться к адресу %v.", address)
	}

	var lock sync.Mutex
	lock.Lock()
	peers.Add(Peer{
		Address: address,
		Client:  client,
	})
	lock.Unlock()
}

func StopServer() {
	for _, peer := range peers.Peers {
		if peer.Client == nil {
			continue
		}
		response, err := peer.Client.Disconnect(context.Background(), &pb.DisconnectRequest{
			Address: myAddress,
		})
		if err == nil && response.Address == peer.Address {
			log.Printf("DISCONNECT: Отсоединились от пира %v.", peer.Address)
		} else {
			log.Printf("ERROR: Не смогли сообщить об отсоединении пиру %v", peer.Address)
		}
	}
}

func StartServer() {
	conn, err := grpc.Dial(config.Config.P2PAddress+":"+strconv.Itoa(config.Config.InitialPort), grpc.WithInsecure())
	if err != nil {
		panic(err)
	}
	client := pb.NewP2PClient(conn)
	response, err := client.List(context.Background(), &pb.ListRequest{})
	if err == nil {
		var addresses []string = response.Addresses
		var lastAddress string = addresses[len(addresses)-1]
		var s []string = strings.Split(lastAddress, ":")
		var _, port string = s[0], s[1]
		intPort, err := strconv.Atoi(port)
		if err != nil {
			panic(err)
		}

		go createServer(intPort + 1)
		for _, address := range addresses {
			go createClient(address)
		}
	} else {
		go createServer(config.Config.InitialPort)
	}
}

func GetInformationFunction(word string) []string {
	allSentences := []string{}
	for _, peer := range peers.Peers {
		if peer.Client == nil {
			continue
		}
		response, err := peer.Client.Get(context.Background(), &pb.GetRequest{
			Address: myAddress,
			Word:    word,
		})
		if err != nil {
			log.Printf("ERROR: не смогли получить информацию о слове %v от пира с адресом %v.", word, peer.Address)
			return []string{}
		}
		allSentences = append(allSentences, response.Sentences...)
	}
	return allSentences
}

var GetInformation = GetInformationFunction
