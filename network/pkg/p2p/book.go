package p2p

import (
	"context"
	"log"

	"github.com/izveigor/p2p-words/network/pkg/p2p/pb"
)

type Book struct {
	Name            string
	CharactersCount int
	WordsCount      int
	Tree            *Tree
}

type Library struct {
	Books []*Book
}

func (library *Library) SearchSentences(word string) (sentences []string) {
	partSentences := make(chan []string)
	for _, book := range library.Books {
		go book.SearchSentences(word, partSentences)
	}

	for i := 0; i < len(library.Books); i++ {
		getSentences := <-partSentences
		sentences = append(sentences, getSentences...)
	}

	return sentences
}

func (library *Library) GetBooksInformation() *pb.GetBooksInformationResponse {
	books := []*pb.Book{}
	for _, book := range library.Books {
		books = append(books, &pb.Book{
			Name:            book.Name,
			CharactersCount: int32(book.CharactersCount),
			WordsCount:      int32(book.WordsCount),
		})
	}

	return &pb.GetBooksInformationResponse{
		Number: int32(len(library.Books)),
		Books:  books,
	}
}

func (book *Book) SearchSentences(word string, partSentences chan<- []string) {
	element := book.Tree.Search(word)

	if element == nil {
		partSentences <- []string{}
	} else {
		partSentences <- element.Sentences
	}
}

func CreateNewBookFunction(ServiceClient Client, tree *Tree, name, text string) (*Book, []*TreeElement) {
	response, err := ServiceClient.Service.Lemmatize(context.Background(), &pb.LemmatizerRequest{
		Text: text,
	})

	if err != nil {
		log.Fatalf("Книга с названием %q не была загружена: ошибка с лемматизацией: %v.", name, err)
	}

	elements, charactersCount, wordsCount := CreateTreeElements(response)

	tree.Create(elements, nil, true)

	return &Book{
		Name:            name,
		CharactersCount: charactersCount,
		WordsCount:      wordsCount,
		Tree:            tree,
	}, elements
}

var CreateNewBook = CreateNewBookFunction
var UserLibrary = Library{[]*Book{}}
