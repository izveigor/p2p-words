package p2p

import (
	"context"
	"github.com/izveigor/p2p-words/network/pkg/config"
	"github.com/izveigor/p2p-words/network/pkg/p2p/pb"
	"github.com/stretchr/testify/assert"

	"google.golang.org/grpc"
	"testing"
	"time"
)

func TestCreateBook(t *testing.T) {
	go InitHTTPServiceServer(config.Config.P2PSvcUrl)
	time.Sleep(time.Millisecond * 1250)

	conn, err := grpc.Dial(config.Config.P2PSvcUrl, grpc.WithInsecure())

	if err != nil {
		t.Fatal(err)
	}

	defer func() {
		CreateFiles = CreateFilesFunction
		UserLibrary = Library{[]*Book{}}
		conn.Close()
	}()

	CreateFiles = func(name string, text []byte) {
		assert.Equal(t, name, "Название")
		assert.Equal(t, text, []byte("Текст"))
	}

	book := &Book{
		Name:            "Название",
		CharactersCount: 50,
		WordsCount:      20,
		Tree:            &Tree{nil},
	}
	UserLibrary = Library{[]*Book{}}
	UserLibrary.Books = append(UserLibrary.Books, book)

	client := pb.NewHTTPClient(conn)
	bookInformation, err := client.CreateBook(context.Background(), &pb.CreateBookRequest{
		Text: []byte("Текст"),
		Name: "Название",
	})

	if err != nil {
		t.Fatal(err)
	}
	assert.Equal(t, bookInformation.Ok, true)
}

func TestGetBooksInformation(t *testing.T) {
	var address string = "localhost:50053"
	go InitHTTPServiceServer(address)
	time.Sleep(time.Millisecond * 1250)

	conn, err := grpc.Dial(address, grpc.WithInsecure())

	if err != nil {
		t.Fatal(err)
	}

	defer func() {
		UserLibrary = Library{[]*Book{}}
		conn.Close()
	}()

	UserLibrary = Library{[]*Book{}}
	UserLibrary.Books = append(UserLibrary.Books, &Book{
		Name:            "Название",
		CharactersCount: 50,
		WordsCount:      20,
		Tree:            &Tree{nil},
	})

	client := pb.NewHTTPClient(conn)
	booksInformation, err := client.GetBooksInformation(context.Background(), &pb.GetBooksInformationRequest{})
	if err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, booksInformation.Number, int32(1))
	assert.Equal(t, booksInformation.Books, []*pb.Book{
		{
			Name:            "Название",
			CharactersCount: 50,
			WordsCount:      20,
		},
	})
}

func TestSearchSentences(t *testing.T) {
	var address string = "localhost:50054"
	go InitHTTPServiceServer(address)
	time.Sleep(time.Millisecond * 1250)

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

	conn, err := grpc.Dial(address, grpc.WithInsecure())

	if err != nil {
		t.Fatal(err)
	}

	defer func() {
		UserLibrary = Library{[]*Book{}}
		GetInformation = GetInformationFunction
		conn.Close()
	}()
	GetInformation = func(word string) []string {
		return []string{"Другое предложение."}
	}

	UserLibrary = Library{[]*Book{}}
	UserLibrary.Books = append(UserLibrary.Books, firstBook)

	client := pb.NewHTTPClient(conn)
	response, err := client.SearchSentences(context.Background(), &pb.SearchSentencesRequest{
		Word: "слово",
	})
	if err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, response.Sentences, []string{
		"Другое предложение.",
		"Предложение",
	})
}
