proto:
	protoc -I ./pkg/p2p/pb/ --go_out=. --go-grpc_out=. ./pkg/p2p/pb/*.proto

mock:
	mockgen -source=./pkg/p2p/pb/lemmatizer_grpc.pb.go -destination=./pkg/mock/lemmatizer_grpc_mock.go -package=mock

start:
	go run cmd/main.go

test:
	go test -v -cover github.com/izveigor/p2p-words/network/pkg/p2p