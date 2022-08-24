proto:
	protoc -I ./pkg/server/pb/ --go_out=. --go-grpc_out=. ./pkg/server/pb/*.proto

mock:
	mockgen -source=./pkg/server/pb/http_grpc.pb.go -destination=./pkg/mock/http_grpc_mock.go -package=mock

start:
	go run cmd/main.go

test:
	go test -v -cover github.com/izveigor/p2p-words/http/pkg/server